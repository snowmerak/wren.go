package main

import (
	"context"
	"fmt"
	"time"

	wren "github.com/snowmerak/wren.go"
)

func main() {
	fmt.Println("=== Async Examples ===\n")
	
	// Example 1: Basic async task
	basicAsyncExample()
	
	// Example 2: Multiple concurrent tasks
	concurrentTasksExample()
	
	// Example 3: Task with timeout
	timeoutExample()
	
	// Example 4: Task cancellation
	cancellationExample()
	
	// Example 5: Waiting for all tasks
	waitAllExample()
	
	fmt.Println("\n=== All async examples completed! ===")
}

func basicAsyncExample() {
	fmt.Println("1. Basic Async Task:")
	
	am := wren.GetAsyncManager()
	
	future := am.Submit(func(ctx context.Context) (interface{}, error) {
		fmt.Println("  Task started...")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("  Task completed!")
		return "Hello from async task!", nil
	})
	
	fmt.Printf("  Future ID: %d\n", future.ID())
	fmt.Printf("  Is ready? %v\n", future.IsReady())
	
	result, err := future.Wait()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	
	fmt.Printf("  Result: %v\n", result)
	fmt.Printf("  Is ready? %v\n\n", future.IsReady())
}

func concurrentTasksExample() {
	fmt.Println("2. Multiple Concurrent Tasks:")
	
	am := wren.GetAsyncManager()
	
	futures := make([]*wren.Future, 5)
	for i := 0; i < 5; i++ {
		taskNum := i + 1
		futures[i] = am.Submit(func(ctx context.Context) (interface{}, error) {
			delay := time.Duration(taskNum*100) * time.Millisecond
			time.Sleep(delay)
			return fmt.Sprintf("Task %d done", taskNum), nil
		})
		fmt.Printf("  Submitted task %d (ID: %d)\n", taskNum, futures[i].ID())
	}
	
	fmt.Println("  Waiting for all tasks...")
	for i, future := range futures {
		result, err := future.Wait()
		if err != nil {
			fmt.Printf("  Task %d error: %v\n", i+1, err)
		} else {
			fmt.Printf("  %v\n", result)
		}
	}
	fmt.Println()
}

func timeoutExample() {
	fmt.Println("3. Task with Timeout:")
	
	am := wren.GetAsyncManager()
	
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	
	future := am.SubmitWithContext(ctx, func(ctx context.Context) (interface{}, error) {
		select {
		case <-time.After(1 * time.Second):
			return "Should not reach here", nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})
	
	result, err := future.Wait()
	if err != nil {
		fmt.Printf("  Task failed as expected: %v\n", err)
	} else {
		fmt.Printf("  Unexpected result: %v\n", result)
	}
	fmt.Println()
}

func cancellationExample() {
	fmt.Println("4. Task Cancellation:")
	
	am := wren.GetAsyncManager()
	
	future := am.Submit(func(ctx context.Context) (interface{}, error) {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				fmt.Printf("  Working... %d/10\n", i+1)
			}
		}
		return "Completed", nil
	})
	
	// Cancel after a short delay
	time.Sleep(250 * time.Millisecond)
	fmt.Println("  Cancelling task...")
	future.Cancel()
	
	result, err := future.Wait()
	if err != nil {
		fmt.Printf("  Task was cancelled: %v\n", err)
	} else {
		fmt.Printf("  Unexpected result: %v\n", result)
	}
	fmt.Println()
}

func waitAllExample() {
	fmt.Println("5. Wait for All Tasks:")
	
	am := wren.NewAsyncManager(2) // Only 2 workers
	defer am.Shutdown()
	
	start := time.Now()
	
	for i := 0; i < 6; i++ {
		taskNum := i + 1
		am.Submit(func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("  Task %d completed\n", taskNum)
			return taskNum, nil
		})
	}
	
	fmt.Println("  Waiting for all tasks to complete...")
	am.WaitAll()
	
	duration := time.Since(start)
	fmt.Printf("  All tasks completed in %v\n", duration)
	fmt.Println()
}
