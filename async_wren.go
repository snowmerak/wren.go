package wrengo

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// AsyncForeignMethodFn is a foreign method function that returns a Future.
type AsyncForeignMethodFn func(vm *WrenVM) (*Future, error)

// asyncForeignMethods stores registered async foreign methods.
var asyncForeignMethods = make(map[string]AsyncForeignMethodFn)
var asyncMethodMutex sync.RWMutex

// RegisterAsyncForeignMethod registers an async foreign method.
func RegisterAsyncForeignMethod(module, className string, isStatic bool, signature string, fn AsyncForeignMethodFn) {
	asyncMethodMutex.Lock()
	defer asyncMethodMutex.Unlock()

	// Create full signature with static prefix if needed
	fullSig := signature
	if isStatic {
		fullSig = "static " + signature
	}

	key := module + ":" + className + ":" + fullSig
	asyncForeignMethods[key] = fn
}

// lookupAsyncForeignMethod finds a registered async foreign method.
func lookupAsyncForeignMethod(module, className string, isStatic bool, signature string) AsyncForeignMethodFn {
	asyncMethodMutex.RLock()
	defer asyncMethodMutex.RUnlock()

	fullSig := signature
	if isStatic {
		fullSig = "static " + signature
	}

	key := module + ":" + className + ":" + fullSig
	return asyncForeignMethods[key]
}

// Wren-facing async API implementation

// asyncRun implements Async.run(_) - executes a function asynchronously.
//
//wren:bind module=main class=Async name=run(_) static
func asyncRun(vm *WrenVM) error {
	// Get the function from slot 1
	// This is a Wren function/closure that we need to call asynchronously

	// For now, we'll create a future that the Wren code can poll
	// In a more complete implementation, we'd store the Wren Fiber/Function
	// and resume it when the async work is done

	futureID := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
		// This is where we would call back into Wren
		// For now, return a placeholder
		return "async_result", nil
	}).ID()

	// Return the future ID to Wren
	vm.SetSlotDouble(0, float64(futureID))
	return nil
}

// asyncCall implements Async.call(name, args) - calls a registered async Go function.
//
//wren:bind module=main class=Async name=call(_,_) static
func asyncCall(vm *WrenVM) error {
	// Get method name from slot 1
	methodName := vm.GetSlotString(1)

	// Get args list from slot 2
	// For simplicity, we'll support basic types
	var args []interface{}
	if vm.GetSlotType(2) == TypeList {
		count := vm.GetListCount(2)
		args = make([]interface{}, count)
		for i := 0; i < count; i++ {
			vm.GetListElement(2, i, 3)

			switch vm.GetSlotType(3) {
			case TypeNum:
				args[i] = vm.GetSlotDouble(3)
			case TypeString:
				args[i] = vm.GetSlotString(3)
			case TypeBool:
				args[i] = vm.GetSlotBool(3)
			}
		}
	}

	// Look up the async method
	asyncMethodMutex.RLock()
	asyncFn := asyncForeignMethods[methodName]
	asyncMethodMutex.RUnlock()

	if asyncFn == nil {
		return fmt.Errorf("async method not found: %s", methodName)
	}

	// Execute asynchronously
	future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
		return asyncFn(vm)
	})

	// Return future ID
	vm.SetSlotDouble(0, float64(future.ID()))
	return nil
}

// asyncAwait implements Async.await(futureId) - waits for a future to complete.
//
//wren:bind module=main class=Async name=await(_) static
func asyncAwait(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}

	result, err := future.Wait()
	if err != nil {
		return err
	}

	// Set result in slot 0
	return setSlotValue(vm, 0, result)
}

// asyncIsReady implements Async.isReady(futureId) - checks if future is ready.
//
//wren:bind module=main class=Async name=isReady(_) static
func asyncIsReady(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}

	vm.SetSlotBool(0, future.IsReady())
	return nil
}

// asyncGet implements Async.get(futureId) - gets result without waiting.
//
//wren:bind module=main class=Async name=get(_) static
func asyncGet(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}

	result, err := future.Get()
	if err != nil {
		return err
	}

	return setSlotValue(vm, 0, result)
}

// asyncCancel implements Async.cancel(futureId) - cancels a future.
//
//wren:bind module=main class=Async name=cancel(_) static
func asyncCancel(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}

	future.Cancel()
	return nil
}

// asyncGetState implements Async.getState(futureId) - gets future state.
//
//wren:bind module=main class=Async name=getState(_) static
func asyncGetState(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}

	state := future.State()
	vm.SetSlotDouble(0, float64(state))
	return nil
}

// asyncCleanup implements Async.cleanup(futureId) - removes a completed future.
//
//wren:bind module=main class=Async name=cleanup(_) static
func asyncCleanup(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)

	GetAsyncManager().RemoveFuture(int64(futureID))
	return nil
}

// Helper function to set slot value based on type
func setSlotValue(vm *WrenVM, slot int, value interface{}) error {
	if value == nil {
		vm.SetSlotNull(slot)
		return nil
	}

	switch v := value.(type) {
	case bool:
		vm.SetSlotBool(slot, v)
	case int:
		vm.SetSlotDouble(slot, float64(v))
	case int64:
		vm.SetSlotDouble(slot, float64(v))
	case float64:
		vm.SetSlotDouble(slot, v)
	case string:
		vm.SetSlotString(slot, v)
	case error:
		vm.SetSlotString(slot, v.Error())
	default:
		vm.SetSlotString(slot, fmt.Sprintf("%v", v))
	}

	return nil
}

// RegisterAsyncMethod registers a simple async method that can be called from Wren.
func RegisterAsyncMethod(name string, fn func(args ...interface{}) (interface{}, error)) {
	asyncMethodMutex.Lock()
	defer asyncMethodMutex.Unlock()

	asyncForeignMethods[name] = func(vm *WrenVM) (*Future, error) {
		// This is a simplified wrapper - in practice you'd extract args from VM
		future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
			return fn()
		})
		return future, nil
	}
}
