package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // ✅ Worker bolta hai: "Mera kaam ho gaya"

	for i := 1; i <= 3; i++ {
		fmt.Printf("Worker %d → step %d\n", id, i)
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup // ✅ Register banaya

	wg.Add(3) // ✅ Bataya ki 3 workers aane wale hain

	go worker(1, &wg)
	go worker(2, &wg)
	go worker(3, &wg)

	wg.Wait() // ✅ Ruk jao jab tak sab Done na bol dein

	fmt.Println("✅ All workers finished")
}