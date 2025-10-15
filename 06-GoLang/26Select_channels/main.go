// /*
//    `select` statement Go language me ek special control structure hai jo multiple channels ko ek saath monitor karta hai. 
//    Jis channel se sabse pehle data aata hai, `select` usi case ko execute karta hai.
// */

// package main

// import (
// 	"fmt"
// 	"time"
// )

// // channel 1 se data bhejne wala function
// func sendToCh1(ch chan string) {
// 	time.Sleep(2 * time.Second)
// 	ch <- "Message from channel 1"
// }

// // channel 2 se data bhejne wala function
// func sendToCh2(ch chan string) {
// 	time.Sleep(1 * time.Second)
// 	ch <- "Message from channel 2"
// }

// func main() {
// 	ch1 := make(chan string)
// 	ch2 := make(chan string)

// 	// go routines start
// 	go sendToCh1(ch1)
// 	go sendToCh2(ch2)

// 	time.Sleep(2 * time.Second) // 👈 Give goroutines enough time to send!


// 	// select block
// 	select {
// 	case msg1 := <-ch1:
// 		fmt.Println(msg1)
// 	case msg2 := <-ch2:
// 		fmt.Println(msg2)
// 	default:
// 		fmt.Println("No message yet")
// 	}
// }


package main

import (
	"fmt"
	"time"
)

func chan1(channel1 chan string){
  
	channel1 <- "Rahman bhai"
	time.Sleep(time.Second*1)

}

func chan2(channel2 chan string){

	channel2 <- "Boss Rahman"
	time.Sleep(time.Second*3)

}


func main(){
 
	channel1 := make(chan string)
	go chan1(channel1)
	 
	channel2 := make(chan string)
	go chan2(channel2)

	time.Sleep(time.Second *2)
	select {
	case msg1 := <- channel1 : 
	   fmt.Println(msg1)
	case msg2 := <- channel2 :
       fmt.Println(msg2)
	default :
	fmt.Println("No Channel came bro")   
	}
}