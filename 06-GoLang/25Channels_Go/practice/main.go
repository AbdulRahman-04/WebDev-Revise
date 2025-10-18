// // package main

// // import (
// // 	"fmt"
// // 	"time"
// // )

// // func add(a, b int)  {
// // 	fmt.Println("addition is going on")
// // 	time.Sleep(500*time.Millisecond)
// // 	fmt.Println(a+b)
// // }

// // func sub(a, b int)  {
// // 	fmt.Println("subtraction is going on")
// // 	time.Sleep(500*time.Millisecond)
// //     fmt.Println(a-b)

// // }

// // func main(){

// // 	go add(150, 478)
// // 	go sub(298, 113)

// // 	time.Sleep(2 *time.Second)
// // }

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func worker1(id int, wg *sync.WaitGroup){
// 	defer wg.Done()

// 	fmt.Printf("worker %d is on\n", id)
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Printf("worker %d is finsihed\n", id)
// }

// func worker2(id int, wg *sync.WaitGroup){
// 	defer wg.Done()

// 	fmt.Printf("worker %d is on\n", id)
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Printf("worker %d is finsihed\n", id)
// }

// func worker3(id int, wg *sync.WaitGroup){
// 	defer wg.Done()
// 	fmt.Printf("worker %d is on\n", id)
// 	time.Sleep(500*time.Millisecond)
// 	fmt.Printf("worker %d is finsihed\n", id)
// }

// func main(){

// 	var wg sync.WaitGroup

// 	wg.Add(3)

// 	go worker1(1, &wg)
// 	go worker2(2,&wg)
// 	go worker3(3,&wg)

// 	wg.Wait()
// }

package main

import "fmt"

func read(ch chan int){
  
	for value := range ch{
		fmt.Println(value)
	}
	fmt.Println("error reading channel")

}

func main(){
	ch := make(chan int)

	go read(ch)

	ch <- 46
    ch <- 49
	ch <- 42
	ch <- 40	
}