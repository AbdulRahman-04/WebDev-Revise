package main

import "fmt"

func main(){
	name := "raheem"
	age := 21
	isAlive := true
	// clg := "dcet"
	marks := 78.65

	// fmt.Println("name", name)
	// fmt.Println("age:", age)
	// fmt.Println("isAlive:", isAlive)
	// fmt.Println("clg:", clg)
	// fmt.Println("marks:", marks)


    fmt.Println("%s", age , isAlive)
	fmt.Println(marks)

	fmt.Printf("name is: %s\n", name)
	fmt.Printf("age is %d:\n", age)
	fmt.Printf("isAlive : %t\n", isAlive)
    fmt.Printf("marks is: %f\n", marks)

}