package main

import (
	"fmt"
	"time"
)

func main(){
 
	timeNow := time.Now()

	formatted := timeNow.Format("Monday 02-Jan-2006 15:04:06")
	fmt.Println(formatted)

}