package mySelect

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


func MySelect(){
 
	channel1 := make(chan string)
	go chan1(channel1)
	 
	channel2 := make(chan string)
	go chan2(channel2)

	time.Sleep(time.Second *3)
	select {
	case msg1 := <- channel1 : 
	   fmt.Println(msg1)
	case msg2 := <- channel2 :
       fmt.Println(msg2)
	default :
	fmt.Println("No Channel came bro")   
	}
}