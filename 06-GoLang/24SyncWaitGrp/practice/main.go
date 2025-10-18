package main

import (
	"fmt"
	"sync"
	"time"
)

func w1(id int, wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Printf("Worker %d start\n", id) // Kaam start
	time.Sleep(1 * time.Second)         // Kaam simulate
	fmt.Printf("Worker %d done\n", id)

}
func w2(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d start\n", id) // Kaam start
	time.Sleep(1 * time.Second)         // Kaam simulate
	fmt.Printf("Worker %d done\n", id)
}

func w3(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d start\n", id) // Kaam start
	time.Sleep(1 * time.Second)         // Kaam simulate
	fmt.Printf("Worker %d done\n", id)
}
func w4(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d start\n", id) // Kaam start
	time.Sleep(1 * time.Second)         // Kaam simulate
	fmt.Printf("Worker %d done\n", id)
}

func main() {

	var wg sync.WaitGroup

	wg.Add(4)

	go w1(1, &wg)
	go w2(2, &wg)
	go w3(3, &wg)
	go w4(4, &wg)

	wg.Wait()

}
