package main

import (
	"fmt"
	"sync"
	"time"
)

// Step 4: Worker function
// value nd mem address recieve kro
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // ✅ Kaam khatam hone pe counter -1 kar do

	fmt.Printf("Worker %d start\n", id)     // Kaam start
	time.Sleep(1 * time.Second)             // Kaam simulate
	fmt.Printf("Worker %d done\n", id)      // Kaam complete
}

func main() {
	// Step 1: WaitGroup variable banaye
	var wg sync.WaitGroup // Tracker for goroutines

	// Step 2: Bataye kitne goroutines track karne hain
	wg.Add(3) // 3 goroutines ka counter set

	// Step 3: Goroutines start karo aur WaitGroup ka address pass karo
	go worker(1, &wg)
	go worker(2, &wg)
	go worker(3, &wg)

	// Step 5: Main goroutine wait kare jab tak sab goroutines finish na ho jaye
	wg.Wait() // Ruk jao yaha jab tak counter 0 na ho

	// Step 6: Sab complete hone ke baad aage badho
	fmt.Println("All workers finished ✅")
}
