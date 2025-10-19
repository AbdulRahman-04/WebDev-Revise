// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func work1(a, b int, wg*sync.WaitGroup){
//   defer wg.Done()
// 	fmt.Println("Worker 1 started")
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Println(a+b)

// }
// func work2(c, d int, wg *sync.WaitGroup){
// 	defer wg.Done()
// fmt.Println("Worker 2 started")
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Println(c-d)
// }
// func work3(x, y int, wg*sync.WaitGroup){
//   defer wg.Done()
// 	fmt.Println("Worker 3 started")
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Println(x*y)

// }

// func main(){

// 	var wg sync.WaitGroup

// 	wg.Add(3)

// 	go work1(12, 23, &wg)
// 	go work2(45, 45, &wg)
// 	go work3(12, 15, &wg)

// 	wg.Wait()

// }

// package main

// import "fmt"

// func readData(ch chan int){
//   for value := range ch {
// 	fmt.Println(value)
//   }
// }

// func main(){
//   ch := make(chan int)

//   go readData(ch)

//   ch <- 5049
//   ch <- 7565
//   ch <- 5467

//   close(ch)
// }

package main

import (
	"fmt"
	"sync"
)

func readData(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range ch {
		fmt.Println(value)
	}
}

func main() {

	var wg sync.WaitGroup

	wg.Add(1)

	ch := make(chan string, 4)

	go readData(ch, &wg)

	ch <- "hey"
	ch <- "bro"
	ch <- "kaise"
	ch <- "ho?"

	close(ch)
	wg.Wait()

	

	

}
