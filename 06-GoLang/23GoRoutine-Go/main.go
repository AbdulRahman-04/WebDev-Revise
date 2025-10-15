package main

import (
	"fmt"
	"time"
)

func greet() {
	for i := 1; i <= 5; i++ {
		fmt.Println("👋 Hello from goroutine", i)
		time.Sleep(500 * time.Millisecond)
	}
}

func work() {
	for i := 1; i <= 5; i++ {
		fmt.Println("💼 Working...", i)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	go greet()
	go work()

	fmt.Println("💬 Main function running...")
	time.Sleep(4 * time.Second)
	fmt.Println("✅ Main function ended")
}


