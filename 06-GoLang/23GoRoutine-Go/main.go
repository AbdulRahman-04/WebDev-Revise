package main

import (
	"fmt"
	"time"
)

func Task(name string){
	for i := 0; i <= 3; i++ {
		 fmt.Println(name)
		 time.Sleep(time.Millisecond*200)
	}
}

func main(){
	go Task("do golang")
	go Task("do go gin mongo project")
	go Task("do frontend")

	fmt.Println("code finished")

	time.Sleep(time.Second*1)
}