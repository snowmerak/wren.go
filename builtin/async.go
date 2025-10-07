package builtin

//go:generate go run ../wrengen -dir .
//go:generate ./wrenlsp-gen.exe ../gwen/builtin_wren.go

import (
	"context"
	"fmt"
	"time"

	wrengo "github.com/snowmerak/wren.go"
)

// Async provides built-in async utilities
//
//wren:bind module=async
type Async struct{}

// Sleep creates a future that completes after specified seconds
//
//wren:bind name=sleep(_) static
func (a *Async) Sleep(vm *wrengo.WrenVM) error {
	seconds := vm.GetSlotDouble(1)

	if seconds < 0 {
		return fmt.Errorf("sleep duration cannot be negative: %f", seconds)
	}

	future := wrengo.GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
		duration := time.Duration(seconds * float64(time.Second))

		// Context 취소 확인과 함께 sleep
		select {
		case <-time.After(duration):
			return fmt.Sprintf("Slept for %.2f seconds", seconds), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})

	vm.SetSlotDouble(0, float64(future.ID()))
	return nil
}

// Delay creates a future that completes after specified milliseconds
//
//wren:bind name=delay(_) static
func (a *Async) Delay(vm *wrengo.WrenVM) error {
	milliseconds := vm.GetSlotDouble(1)

	if milliseconds < 0 {
		return fmt.Errorf("delay duration cannot be negative: %f", milliseconds)
	}

	future := wrengo.GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
		duration := time.Duration(milliseconds * float64(time.Millisecond))

		select {
		case <-time.After(duration):
			return fmt.Sprintf("Delayed for %.0f milliseconds", milliseconds), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})

	vm.SetSlotDouble(0, float64(future.ID()))
	return nil
}

// Timer creates a future that acts as a simple timer
//
//wren:bind name=timer(_,_) static
func (a *Async) Timer(vm *wrengo.WrenVM) error {
	duration := vm.GetSlotDouble(1) // seconds
	message := vm.GetSlotString(2)  // custom message

	if duration < 0 {
		return fmt.Errorf("timer duration cannot be negative: %f", duration)
	}

	future := wrengo.GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
		d := time.Duration(duration * float64(time.Second))

		select {
		case <-time.After(d):
			if message == "" {
				return fmt.Sprintf("Timer completed after %.2f seconds", duration), nil
			}
			return message, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})

	vm.SetSlotDouble(0, float64(future.ID()))
	return nil
}
