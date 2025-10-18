package main

import (
	"fmt"
	"time"
)

func Work(){
	for i:=0; i<= 10; i++ {
		fmt.Println("Work1 is done")

		time.Sleep(1000*time.Millisecond)
	}
}

func Work2(){
	for j := 0; j <= 10; j ++ {
		fmt.Println("WORK 2 is done")

		time.Sleep(1000*time.Millisecond)
	}
}

func main(){
	go Work()
	go Work2()


	time.Sleep(15*time.Second)
}