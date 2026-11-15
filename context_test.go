package belajargolangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextC, "e", "E")
	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))
	fmt.Println(contextA.Value("b"))
}

func CreateCounterLeak() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()
	return destination
}
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) //simulasi slow query
			}
		}
	}()
	return destination
}

func TestGoroutineLeak(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	destination := CreateCounterLeak()
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}

	fmt.Println(runtime.NumGoroutine())
}
func TestGoroutineWithCancel(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)
	destination := CreateCounter(ctx)
	fmt.Println(runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 {
			break
		}
	}
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}
func TestGoroutineWithTimeOut(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println(runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("counter", n)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}
func TestGoroutineWithDeadline(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second)) // bisa di setting sesuka hati, misal di jam 00.00
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println(runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("counter", n)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}
