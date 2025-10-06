package wrengo

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// FutureState represents the state of a Future.
type FutureState int

const (
	// FuturePending indicates the future is still running.
	FuturePending FutureState = iota
	// FutureCompleted indicates the future has completed successfully.
	FutureCompleted
	// FutureFailed indicates the future has failed with an error.
	FutureFailed
	// FutureCancelled indicates the future was cancelled.
	FutureCancelled
)

// Future represents an asynchronous computation result.
type Future struct {
	id        int64
	state     atomic.Int32
	result    atomic.Value
	err       atomic.Value
	done      chan struct{}
	ctx       context.Context
	cancel    context.CancelFunc
	startedAt int64
}

// newFuture creates a new Future with a unique ID.
func newFuture(ctx context.Context) *Future {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)

	f := &Future{
		id:     atomic.AddInt64(&nextFutureID, 1),
		done:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
	}
	f.state.Store(int32(FuturePending))

	return f
}

// ID returns the unique ID of this future.
func (f *Future) ID() int64 {
	return f.id
}

// State returns the current state of the future.
func (f *Future) State() FutureState {
	return FutureState(f.state.Load())
}

// IsReady returns true if the future has completed (successfully or with error).
func (f *Future) IsReady() bool {
	state := f.State()
	return state != FuturePending
}

// Wait blocks until the future completes and returns the result or error.
func (f *Future) Wait() (interface{}, error) {
	<-f.done

	state := f.State()
	switch state {
	case FutureCompleted:
		return f.result.Load(), nil
	case FutureFailed:
		if err := f.err.Load(); err != nil {
			return nil, err.(error)
		}
		return nil, errors.New("future failed with unknown error")
	case FutureCancelled:
		return nil, errors.New("future was cancelled")
	default:
		return nil, errors.New("future in invalid state")
	}
}

// Get returns the result if ready, or an error if not ready or failed.
func (f *Future) Get() (interface{}, error) {
	if !f.IsReady() {
		return nil, errors.New("future not ready")
	}

	state := f.State()
	switch state {
	case FutureCompleted:
		return f.result.Load(), nil
	case FutureFailed:
		if err := f.err.Load(); err != nil {
			return nil, err.(error)
		}
		return nil, errors.New("future failed with unknown error")
	case FutureCancelled:
		return nil, errors.New("future was cancelled")
	default:
		return nil, errors.New("future in invalid state")
	}
}

// Cancel cancels the future's context.
func (f *Future) Cancel() {
	if f.IsReady() {
		return
	}
	f.cancel()
	f.setState(FutureCancelled)
	close(f.done)
}

// Context returns the future's context.
func (f *Future) Context() context.Context {
	return f.ctx
}

// complete marks the future as completed with a result.
func (f *Future) complete(result interface{}) {
	if f.IsReady() {
		return
	}
	f.result.Store(result)
	f.setState(FutureCompleted)
	close(f.done)
}

// fail marks the future as failed with an error.
func (f *Future) fail(err error) {
	if f.IsReady() {
		return
	}
	if err == nil {
		err = errors.New("unknown error")
	}
	f.err.Store(err)
	f.setState(FutureFailed)
	close(f.done)
}

// setState atomically sets the future's state.
func (f *Future) setState(state FutureState) {
	f.state.Store(int32(state))
}

// AsyncTask represents a task that can be executed asynchronously.
type AsyncTask func(ctx context.Context) (interface{}, error)

// AsyncManager manages asynchronous tasks and futures.
type AsyncManager struct {
	futures sync.Map // map[int64]*Future
	workers int
	queue   chan *asyncJob
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// asyncJob represents a job in the work queue.
type asyncJob struct {
	future *Future
	task   AsyncTask
}

var (
	nextFutureID int64
	globalAsync  *AsyncManager
	asyncOnce    sync.Once
)

// GetAsyncManager returns the global async manager, creating it if necessary.
func GetAsyncManager() *AsyncManager {
	asyncOnce.Do(func() {
		globalAsync = NewAsyncManager(0) // 0 = number of CPU cores
	})
	return globalAsync
}

// NewAsyncManager creates a new async manager with the specified number of workers.
// If workers is 0 or negative, it defaults to runtime.NumCPU().
func NewAsyncManager(workers int) *AsyncManager {
	if workers <= 0 {
		workers = 4 // Default to 4 workers
	}

	ctx, cancel := context.WithCancel(context.Background())

	am := &AsyncManager{
		workers: workers,
		queue:   make(chan *asyncJob, workers*10),
		ctx:     ctx,
		cancel:  cancel,
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		am.wg.Add(1)
		go am.worker()
	}

	return am
}

// worker processes jobs from the queue.
func (am *AsyncManager) worker() {
	defer am.wg.Done()

	for {
		select {
		case <-am.ctx.Done():
			return
		case job := <-am.queue:
			if job == nil {
				return
			}
			am.executeJob(job)
		}
	}
}

// executeJob executes a single job.
func (am *AsyncManager) executeJob(job *asyncJob) {
	defer func() {
		if r := recover(); r != nil {
			job.future.fail(errors.New("panic in async task"))
		}
	}()

	result, err := job.task(job.future.ctx)
	if err != nil {
		job.future.fail(err)
	} else {
		job.future.complete(result)
	}
}

// Submit submits a task for asynchronous execution and returns a Future.
func (am *AsyncManager) Submit(task AsyncTask) *Future {
	return am.SubmitWithContext(nil, task)
}

// SubmitWithContext submits a task with a custom context.
func (am *AsyncManager) SubmitWithContext(ctx context.Context, task AsyncTask) *Future {
	if ctx == nil {
		ctx = am.ctx
	}

	future := newFuture(ctx)
	am.futures.Store(future.ID(), future)

	job := &asyncJob{
		future: future,
		task:   task,
	}

	select {
	case am.queue <- job:
		// Job queued successfully
	case <-am.ctx.Done():
		future.fail(errors.New("async manager is shutting down"))
	}

	return future
}

// GetFuture retrieves a future by its ID.
func (am *AsyncManager) GetFuture(id int64) (*Future, bool) {
	if val, ok := am.futures.Load(id); ok {
		return val.(*Future), true
	}
	return nil, false
}

// RemoveFuture removes a future from the manager.
func (am *AsyncManager) RemoveFuture(id int64) {
	am.futures.Delete(id)
}

// Shutdown gracefully shuts down the async manager.
func (am *AsyncManager) Shutdown() {
	am.cancel()
	close(am.queue)
	am.wg.Wait()
}

// WaitAll waits for all pending futures to complete.
func (am *AsyncManager) WaitAll() {
	var pending []*Future
	am.futures.Range(func(key, value interface{}) bool {
		f := value.(*Future)
		if !f.IsReady() {
			pending = append(pending, f)
		}
		return true
	})

	for _, f := range pending {
		f.Wait()
	}
}
