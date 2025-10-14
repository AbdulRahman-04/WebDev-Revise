package main

import (
	
	"fmt"
	
)


func main(){
	fmt.Println("enter a number to check if odd or even")

	// var myNum int

	// fmt.Scan(&myNum)

	// if myNum%2 != 0 {
	// 	fmt.Println("number is odd")
	// } else {
	// 	fmt.Println("number is even")
	// }

	var checkNum int

	fmt.Scan(&checkNum)

	if checkNum%2!= 0 {
		fmt.Println("Number is odd")
	} else {
		fmt.Println("number is even")
	}
	
}