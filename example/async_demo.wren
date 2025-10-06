// Async API Example in Wren

import "main" for Async

class AsyncDemo {
  // Example: Check future state
  static checkState(futureId) {
    var state = Async.getState(futureId)
    var stateName = "Unknown"
    
    if (state == 0) stateName = "Pending"
    if (state == 1) stateName = "Completed"  
    if (state == 2) stateName = "Failed"
    if (state == 3) stateName = "Cancelled"
    
    System.print("  Future %(futureId) state: %(stateName)")
    return state
  }
  
  // Example: Poll for result (non-blocking)
  static pollFuture(futureId) {
    var attempts = 0
    while (attempts < 10) {
      if (Async.isReady(futureId)) {
        System.print("  Future is ready after %(attempts + 1) checks!")
        var result = Async.get(futureId)
        System.print("  Result: %(result)")
        Async.cleanup(futureId)
        return result
      }
      System.print("  Check %(attempts + 1): Not ready yet...")
      attempts = attempts + 1
    }
    System.print("  Gave up waiting")
    Async.cleanup(futureId)
    return null
  }
  
  // Example: Block until complete
  static blockingWait(futureId) {
    System.print("  Waiting for future %(futureId)...")
    var result = Async.await(futureId)
    System.print("  Got result: %(result)")
    Async.cleanup(futureId)
    return result
  }
}

// Usage example:
// var futureId = ... // Get future ID from Go
// AsyncDemo.blockingWait(futureId)
// or
// AsyncDemo.pollFuture(futureId)
