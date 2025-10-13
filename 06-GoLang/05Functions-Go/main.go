package main

import "fmt"


func simpleFunc(){
  fmt.Println("hello from simple function")	
}

func add(a int, b int) int {
	return a + b;
}

func sub(x, y int) int {
	return x - y
}

func Mul(c int, d int) int {
	return  c * d;
}


func main(){

	simpleFunc()

	result := add(4, 4)
	fmt.Println("Value of a + b is:", result)

	result1 := sub(5, 3)
	fmt.Println("Value after subtraction is :", result1)

	mulResult := Mul(3, 5)
	fmt.Println(mulResult)
}