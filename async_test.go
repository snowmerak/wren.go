package wrengo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFutureBasic(t *testing.T) {
	future := newFuture(nil)
	
	if future.IsReady() {
		t.Error("New future should not be ready")
	}
	
	go func() {
		time.Sleep(50 * time.Millisecond)
		future.complete("test result")
	}()
	
	result, err := future.Wait()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result != "test result" {
		t.Errorf("Expected 'test result', got %v", result)
	}
	
	if !future.IsReady() {
		t.Error("Future should be ready after completion")
	}
}

func TestFutureError(t *testing.T) {
	future := newFuture(nil)
	
	go func() {
		time.Sleep(50 * time.Millisecond)
		future.fail(errors.New("test error"))
	}()
	
	_, err := future.Wait()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFutureCancel(t *testing.T) {
	future := newFuture(nil)
	
	future.Cancel()
	
	_, err := future.Wait()
	if err == nil {
		t.Error("Expected cancellation error")
	}
	
	if future.State() != FutureCancelled {
		t.Errorf("Expected state FutureCancelled, got %v", future.State())
	}
}

func TestAsyncManager(t *testing.T) {
	am := NewAsyncManager(2)
	defer am.Shutdown()
	
	// Submit a simple task
	future := am.Submit(func(ctx context.Context) (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})
	
	if future.IsReady() {
		t.Error("Future should not be ready immediately")
	}
	
	result, err := future.Wait()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestAsyncManagerMultipleTasks(t *testing.T) {
	am := NewAsyncManager(4)
	defer am.Shutdown()
	
	futures := make([]*Future, 10)
	for i := 0; i < 10; i++ {
		idx := i
		futures[i] = am.Submit(func(ctx context.Context) (interface{}, error) {
			time.Sleep(50 * time.Millisecond)
			return idx * 2, nil
		})
	}
	
	for i, future := range futures {
		result, err := future.Wait()
		if err != nil {
			t.Errorf("Task %d failed: %v", i, err)
		}
		
		expected := i * 2
		if result != expected {
			t.Errorf("Task %d: expected %d, got %v", i, expected, result)
		}
	}
}

func TestAsyncManagerContext(t *testing.T) {
	am := NewAsyncManager(2)
	defer am.Shutdown()
	
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	future := am.SubmitWithContext(ctx, func(ctx context.Context) (interface{}, error) {
		select {
		case <-time.After(200 * time.Millisecond):
			return "completed", nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})
	
	_, err := future.Wait()
	if err == nil {
		t.Error("Expected context timeout error")
	}
}

func TestAsyncManagerWaitAll(t *testing.T) {
	am := NewAsyncManager(2)
	defer am.Shutdown()
	
	for i := 0; i < 5; i++ {
		am.Submit(func(ctx context.Context) (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return "done", nil
		})
	}
	
	start := time.Now()
	am.WaitAll()
	duration := time.Since(start)
	
	// With 2 workers, 5 tasks should take at least 200ms (2 batches)
	if duration < 200*time.Millisecond {
		t.Errorf("WaitAll finished too quickly: %v", duration)
	}
}

func TestAsyncManagerGetFuture(t *testing.T) {
	am := NewAsyncManager(2)
	defer am.Shutdown()
	
	future1 := am.Submit(func(ctx context.Context) (interface{}, error) {
		return "result", nil
	})
	
	future2, ok := am.GetFuture(future1.ID())
	if !ok {
		t.Error("Should be able to retrieve future by ID")
	}
	
	if future1.ID() != future2.ID() {
		t.Error("Retrieved future should have same ID")
	}
	
	_, ok = am.GetFuture(99999)
	if ok {
		t.Error("Should not find non-existent future")
	}
}

func TestAsyncManagerRemoveFuture(t *testing.T) {
	am := NewAsyncManager(2)
	defer am.Shutdown()
	
	future := am.Submit(func(ctx context.Context) (interface{}, error) {
		return "result", nil
	})
	
	future.Wait()
	
	am.RemoveFuture(future.ID())
	
	_, ok := am.GetFuture(future.ID())
	if ok {
		t.Error("Future should be removed")
	}
}

func TestGlobalAsyncManager(t *testing.T) {
	am1 := GetAsyncManager()
	am2 := GetAsyncManager()
	
	if am1 != am2 {
		t.Error("GetAsyncManager should return same instance")
	}
	
	future := am1.Submit(func(ctx context.Context) (interface{}, error) {
		return "test", nil
	})
	
	result, err := future.Wait()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}
}

func TestAsyncFutureGet(t *testing.T) {
	future := newFuture(nil)
	
	// Try to get before ready
	_, err := future.Get()
	if err == nil {
		t.Error("Get should fail when future is not ready")
	}
	
	// Complete the future
	future.complete("ready")
	
	// Now Get should succeed
	result, err := future.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result != "ready" {
		t.Errorf("Expected 'ready', got %v", result)
	}
}
