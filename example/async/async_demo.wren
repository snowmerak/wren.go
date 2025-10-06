// Wren Async API Examples
// This demonstrates how to use async operations from Wren

class AsyncDemo {
  // Example 1: Check if a future is ready (non-blocking)
  static pollFuture(futureId) {
    var attempts = 0
    while (attempts < 10) {
      if (Async.isReady(futureId)) {
        System.print("Future is ready!")
        var result = Async.get(futureId)
        System.print("Result: %(result)")
        Async.cleanup(futureId)
        return result
      }
      System.print("Future not ready, attempt %(attempts + 1)")
      attempts = attempts + 1
      // In real code, you'd do other work here
    }
    System.print("Gave up waiting")
    return null
  }
  
  // Example 2: Block until future completes
  static blockingWait(futureId) {
    System.print("Waiting for future %(futureId)...")
    var result = Async.await(futureId)
    System.print("Got result: %(result)")
    Async.cleanup(futureId)
    return result
  }
  
  // Example 3: Check future state
  static checkState(futureId) {
    var state = Async.getState(futureId)
    var stateName = "Unknown"
    
    if (state == 0) stateName = "Pending"
    if (state == 1) stateName = "Completed"
    if (state == 2) stateName = "Failed"
    if (state == 3) stateName = "Cancelled"
    
    System.print("Future %(futureId) state: %(stateName)")
    return state
  }
  
  // Example 4: Cancel a future
  static cancelDemo(futureId) {
    System.print("Cancelling future %(futureId)")
    Async.cancel(futureId)
    checkState(futureId)
  }
}

// Example usage:
// var futureId = Async.call("someMethod", [arg1, arg2])
// AsyncDemo.blockingWait(futureId)
