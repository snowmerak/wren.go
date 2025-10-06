# Async Examples

This directory demonstrates the asynchronous task execution system in wren.go.

## Overview

The async system allows you to:
- Execute Go functions asynchronously from Wren
- Submit tasks to a managed goroutine pool
- Use Future objects to check completion status
- Cancel long-running tasks
- Set timeouts on operations

## Architecture

### Key Components

1. **Future**: Represents an asynchronous computation result
   - States: Pending, Completed, Failed, Cancelled
   - Methods: `Wait()`, `Get()`, `IsReady()`, `Cancel()`

2. **AsyncManager**: Manages worker goroutines and task queue
   - Configurable worker pool size
   - Task submission and tracking
   - Graceful shutdown support

3. **Wren API**: Foreign methods for async operations
   - `Async.call(name, args)` - Execute registered Go function
   - `Async.await(futureId)` - Block until completion
   - `Async.isReady(futureId)` - Check if ready (non-blocking)
   - `Async.get(futureId)` - Get result if ready
   - `Async.cancel(futureId)` - Cancel execution
   - `Async.getState(futureId)` - Get current state
   - `Async.cleanup(futureId)` - Remove completed future

## Examples

### Go Side

```go
// Get the global async manager
am := wren.GetAsyncManager()

// Submit a task
future := am.Submit(func(ctx context.Context) (interface{}, error) {
    // Do heavy computation
    time.Sleep(1 * time.Second)
    return "result", nil
})

// Wait for completion
result, err := future.Wait()

// Or check without blocking
if future.IsReady() {
    result, err := future.Get()
}

// Or cancel
future.Cancel()
```

### Wren Side

```wren
// Call a registered async Go function
var futureId = Async.call("heavyComputation", [arg1, arg2])

// Option 1: Block until done
var result = Async.await(futureId)

// Option 2: Poll for completion
while (!Async.isReady(futureId)) {
    // Do other work
}
var result = Async.get(futureId)

// Option 3: Check state
var state = Async.getState(futureId)
if (state == 1) { // Completed
    var result = Async.get(futureId)
}

// Always cleanup when done
Async.cleanup(futureId)
```

## Running the Examples

### Go Examples
```bash
cd example/async
go run main.go
```

This will demonstrate:
1. Basic async task execution
2. Multiple concurrent tasks
3. Tasks with timeouts
4. Task cancellation
5. Waiting for all tasks

### Wren Examples

See `async_demo.wren` for Wren-side usage patterns.

## Performance Considerations

- **Worker Pool Size**: Default is 4 workers. Adjust based on your workload:
  ```go
  am := wren.NewAsyncManager(8) // 8 workers
  ```

- **Queue Size**: Default is `workers * 10`. Tasks submitted when queue is full will block.

- **Memory**: Completed futures remain in memory until `cleanup()` is called. Always cleanup when done.

- **Context**: Use contexts for timeouts and cancellation:
  ```go
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  future := am.SubmitWithContext(ctx, task)
  ```

## Integration with wrengen

You can register async methods using the `//wren:async` annotation (future feature):

```go
//wren:async module=main class=MyClass name=compute(_)
func asyncCompute(vm *wren.WrenVM, input int) (int, error) {
    // Heavy computation
    time.Sleep(1 * time.Second)
    return input * 2, nil
}
```

## Limitations

1. **No Fiber Integration**: Currently doesn't integrate with Wren's Fiber system. Future work could allow yielding a Fiber until async work completes.

2. **Simple Type Support**: Only basic types (numbers, strings, booleans) are supported for async method arguments and return values.

3. **No Wren Callbacks**: You cannot pass Wren functions as callbacks to async tasks yet.

## Future Enhancements

- [ ] Fiber.yield() integration for true async/await
- [ ] Promise-style chaining (`.then()`, `.catch()`)
- [ ] Async iterators
- [ ] Channel-based streaming results
- [ ] Better type conversion for complex objects
- [ ] Support for Wren function callbacks
- [ ] Automatic cleanup of old futures
- [ ] Priority queue for tasks
- [ ] Task metrics and monitoring

## Thread Safety

All async operations are thread-safe:
- Multiple VMs can share the same AsyncManager
- Future objects can be accessed from multiple goroutines
- The internal state is protected by atomic operations and mutexes

## Error Handling

Errors in async tasks:
1. Are captured and stored in the Future
2. Are returned when calling `Wait()` or `Get()`
3. Set the Future state to `Failed`
4. Include panic recovery

```go
future := am.Submit(func(ctx context.Context) (interface{}, error) {
    return nil, errors.New("task failed")
})

result, err := future.Wait()
// err will be: "task failed"
```
