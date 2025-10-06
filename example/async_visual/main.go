package main

import (
	"context"
	"fmt"
	"time"

	wrengo "github.com/snowmerak/wren.go"
)

func main() {
	fmt.Println("=== Async Non-Blocking Visualization ===\n")
	
	vm := wrengo.NewVMWithForeign()
	defer vm.Free()
	
	// Wren code that demonstrates non-blocking behavior
	wrenCode := `
class Async {
  foreign static await(futureId)
  foreign static isReady(futureId)
  foreign static get(futureId)
  foreign static cleanup(futureId)
}

class SlowTask {
  foreign static start()
}

System.print("[Wren] Step 1: Calling SlowTask.start()...")
var futureId = SlowTask.start()
System.print("[Wren] Step 2: Got future ID: %(futureId) - the call returned immediately!")
System.print("[Wren] Step 3: Now I can do other work while the background task runs...\n")

// Simulate doing other work
System.print("[Wren] Doing other work:")
var i = 0
while (i < 5) {
  System.print("[Wren]   Task %(i + 1)/5 - checking if async is ready...")
  
  if (Async.isReady(futureId)) {
    System.print("[Wren]   ✓ Async task completed while I was working!\n")
    break
  } else {
    System.print("[Wren]   ✗ Not ready yet, continuing with my work...")
  }
  
  i = i + 1
}

System.print("[Wren] Step 4: Checking final result...")
if (Async.isReady(futureId)) {
  var result = Async.get(futureId)
  System.print("[Wren] ✓ Result: %(result)")
  Async.cleanup(futureId)
} else {
  System.print("[Wren] Still not ready, blocking wait now...")
  var result = Async.await(futureId)
  System.print("[Wren] ✓ Result: %(result)")
  Async.cleanup(futureId)
}

System.print("\n[Wren] All done! Notice:")
System.print("[Wren]   1. SlowTask.start() returned immediately")
System.print("[Wren]   2. Wren executed 5 tasks while waiting")
System.print("[Wren]   3. Background work happened in parallel")
`
	
	// Register the SlowTask.start() method that creates an async task
	wrengo.RegisterForeignMethod("main", "SlowTask", true, "start()", func(vm *wrengo.WrenVM) {
		fmt.Println("[Go] >>> SlowTask.start() foreign method called")
		
		// Submit to async manager
		future := wrengo.GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
			duration := 2 * time.Second
			fmt.Printf("[Go Background] Starting heavy computation (will take %v)...\n", duration)
			time.Sleep(duration)
			fmt.Println("[Go Background] Heavy computation finished!")
			return "Computation result: 42", nil
		})
		
		futureID := future.ID()
		fmt.Printf("[Go] <<< Returning future ID %d to Wren immediately\n\n", futureID)
		
		// Return future ID to Wren
		vm.SetSlotDouble(0, float64(futureID))
	})
	
	// Execute Wren code
	fmt.Println("Starting Wren execution...\n")
	fmt.Println("=" + string(make([]byte, 60)) + "=\n")
	for i := range make([]byte, 60) {
		fmt.Print("=")
		_ = i
	}
	fmt.Println("=\n")
	
	result, err := vm.Interpret("main", wrenCode)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	if result != wrengo.ResultSuccess {
		fmt.Printf("Execution failed: %v\n", result)
		return
	}
	
	fmt.Println("\n" + string(make([]byte, 60)) + "=")
	for i := range make([]byte, 60) {
		fmt.Print("=")
		_ = i
	}
	fmt.Println("=")
	
	fmt.Println("\n=== Analysis ===")
	fmt.Println("✓ Foreign method (SlowTask.start) returned immediately")
	fmt.Println("✓ Wren continued executing without blocking")
	fmt.Println("✓ Background task ran in a separate goroutine")
	fmt.Println("✓ Wren polled for completion with Async.isReady()")
	fmt.Println("\nThis demonstrates TRUE non-blocking async in Wren!")
}
