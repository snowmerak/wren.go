# Async Visual Demo

This example demonstrates **TRUE non-blocking async execution** in Wren.

## What This Proves

The demo shows that:

1. **Foreign method calls return immediately** - `SlowTask.start()` submits work and returns a Future ID instantly
2. **Wren continues executing** - While background work runs, Wren executes other code
3. **Background work runs in parallel** - Go goroutines handle heavy computation
4. **Polling is non-blocking** - `Async.isReady()` checks status without waiting

## How It Works

```
[Timeline]

0ms: Wren calls SlowTask.start()
     ↓
1ms: Go receives call, submits to AsyncManager, returns Future ID
     ↓
2ms: Wren continues (got Future ID)
     ├─> Executes Task 1
     ├─> Executes Task 2  
     ├─> Executes Task 3
     ├─> Executes Task 4
     └─> Executes Task 5
     
Meanwhile in background:
     Go goroutine does 2-second computation
     
2000ms: Background work completes
        ↓
2001ms: Wren checks Async.isReady() or Async.await()
        Gets result: "Computation result: 42"
```

## Key Observations

### ✓ Non-Blocking Foreign Method
```wren
var futureId = SlowTask.start()  // Returns immediately!
System.print("Got ID: %(futureId)")  // Executes right away
```

### ✓ Polling While Working
```wren
while (working) {
  doWork()
  if (Async.isReady(futureId)) {
    // Result is ready!
  }
}
```

### ✓ Background Execution
```go
// Go goroutine runs independently
future := am.Submit(func(ctx context.Context) (interface{}, error) {
    time.Sleep(2 * time.Second)  // Heavy work
    return result, nil
})
```

## Running the Demo

```bash
go run ./example/async_visual
```

## Expected Output Pattern

You'll see output interleaved like this:

```
[Go] >>> SlowTask.start() called
[Go] <<< Returning future ID immediately

[Wren] Got future ID: 1
[Wren] Doing other work...
[Wren]   Task 1/5 - not ready yet
[Wren]   Task 2/5 - not ready yet

[Go Background] Starting computation...

[Wren]   Task 3/5 - not ready yet
[Wren]   Task 4/5 - not ready yet

[Go Background] Computation finished!

[Wren]   Task 5/5 - checking...
[Wren] ✓ Result: Computation result: 42
```

The interleaving proves that:
- Wren didn't block waiting for the computation
- Background work ran in parallel
- Wren polled for results asynchronously

## Contrast with Blocking Approach

**Without Async (Blocking):**
```wren
// This would block for 2 seconds
var result = SlowTask.compute()  // ⏳ Wren waits...
System.print("Finally done: %(result)")
```

**With Async (Non-Blocking):**
```wren
// This returns immediately with a Future ID
var futureId = SlowTask.start()  // ✓ Returns instantly!

// Do other work
doSomethingElse()

// Check result later
var result = Async.await(futureId)
```

## Use Cases

This pattern is perfect for:

- **Heavy computations** - CPU-intensive tasks
- **I/O operations** - Network calls, file operations
- **Database queries** - Long-running queries
- **External API calls** - HTTP requests
- **Batch processing** - Processing large datasets

## Architecture

```
┌─────────────────────────────────────┐
│          Wren Code                  │
│                                     │
│  futureId = Task.start() ────┐    │
│  // continues immediately     │    │
│  Async.isReady(futureId) ────┼───┐│
└───────────────────────────────┼───┼┘
                                │   │
                                ▼   │
┌───────────────────────────────────┼┘
│          Go Foreign Method        │
│                                   │
│  Submit to AsyncManager ──────┐  │
│  Return Future ID immediately │  │
└────────────────────────────┬──┘  │
                             │     │
                             ▼     │
┌────────────────────────────────┐ │
│      AsyncManager (Go)         │ │
│                                │ │
│  Worker Goroutine Pool         │ │
│  ├─ Task 1 (running)          │ │
│  ├─ Task 2 (queued)           │ │
│  └─ Task 3 (completed) ───────┼─┘
└────────────────────────────────┘
```

## Conclusion

This demo proves that wren.go's async system provides:

✓ **True non-blocking execution** - Foreign methods return immediately  
✓ **Parallel background processing** - Work happens in Go goroutines  
✓ **Flexible result retrieval** - Poll or block as needed  
✓ **No Wren VM blocking** - VM continues executing while work runs  

Perfect for building responsive applications that need to handle long-running operations!
