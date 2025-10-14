package main

import "fmt"

func simpleFunc(){
	fmt.Println("Hello simple function")
}

func add(x int, y int) int {
	return  x + y
}


func sub(a,b int) int {
	return  a - b
}


func main(){

	simpleFunc()

	addValue := add(25, 25)
	fmt.Println(addValue)

	subValue := sub(15, 12)
	fmt.Println(subValue)
	
}