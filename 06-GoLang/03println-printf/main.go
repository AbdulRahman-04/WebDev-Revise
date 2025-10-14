package main

import "fmt"

func main(){
	name := "raheem"
	age := 21
	isAlive := true
	clg := "dcet"
	marks := 78.65

	// fmt.Println("name", name)
	// fmt.Println("age:", age)
	// fmt.Println("isAlive:", isAlive)
	// fmt.Println("clg:", clg)
	// fmt.Println("marks:", marks)


    fmt.Println("name is:", name)
	fmt.Println("age is: %d", age)

	fmt.Printf("isAlive value is %t\n:", isAlive)
    fmt.Printf("college is %s:\n", clg)

	fmt.Println("name is : %s", name)
	// fmt.Println("heyya")

	fmt.Printf("name is %s\n", name)
	fmt.Printf("age is %d\n", age)
    fmt.Printf("marks is %f\n", marks)
}