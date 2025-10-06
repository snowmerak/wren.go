package wrengo

//go:generate go run ./wrengen -dir .

import (
	"context"
	"errors"
	"fmt"
)

// AsyncForeignMethodFn is a foreign method function that returns a Future.
type AsyncForeignMethodFn func(vm *WrenVM) (*Future, error)

// asyncForeignMethods stores registered async foreign methods.
var asyncForeignMethods = make(map[string]AsyncForeignMethodFn)

// RegisterAsyncMethod registers a simple async method that can be called from Wren.
func RegisterAsyncMethod(name string, fn func(args ...interface{}) (interface{}, error)) {
	asyncForeignMethods[name] = func(vm *WrenVM) (*Future, error) {
		future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
			return fn()
		})
		return future, nil
	}
}

// Async provides asynchronous task execution from Wren.
//wren:bind module=main
type Async struct{}

// Await waits for a future to complete and returns its result.
//wren:bind name=await(_) static
func (a *Async) Await(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)
	
	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}
	
	result, err := future.Wait()
	if err != nil {
		return err
	}
	
	return setSlotValue(vm, 0, result)
}

// IsReady checks if a future is ready without blocking.
//wren:bind name=isReady(_) static
func (a *Async) IsReady(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)
	
	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}
	
	vm.SetSlotBool(0, future.IsReady())
	return nil
}

// Get retrieves the result of a future if it's ready.
//wren:bind name=get(_) static
func (a *Async) Get(vm *WrenVM) error {
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

// Cancel cancels a future.
//wren:bind name=cancel(_) static
func (a *Async) Cancel(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)
	
	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}
	
	future.Cancel()
	return nil
}

// GetState returns the state of a future.
//wren:bind name=getState(_) static
func (a *Async) GetState(vm *WrenVM) error {
	futureID := vm.GetSlotDouble(1)
	
	future, ok := GetAsyncManager().GetFuture(int64(futureID))
	if !ok {
		return errors.New("future not found")
	}
	
	state := future.State()
	vm.SetSlotDouble(0, float64(state))
	return nil
}

// Cleanup removes a completed future from memory.
//wren:bind name=cleanup(_) static
func (a *Async) Cleanup(vm *WrenVM) error {
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
