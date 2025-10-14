package main

import "fmt"

func SimpleFunc(){
	fmt.Println("All functions executed!🚀")
}

func add(x int, y int) int {
	return x + y;
}

func main(){

	defer SimpleFunc()

	result := add(5, 5)
	fmt.Println(result)


	defer fmt.Println("Hey1")
	defer fmt.Println("Hey2")
   defer fmt.Println("Hey3")


}