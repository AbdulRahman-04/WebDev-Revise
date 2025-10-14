package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
 
	// taking inputs from user

	// two types : fmt.scan && bufio 

	// fmt.Println("Type ur name")
	// var name string

	// fmt.Scan(&name)
	// fmt.Printf("Your name is %s\n", name)


	// bufio 

	// fmt.Println("Enter ur username")

	// reader := bufio.NewReader(os.Stdin)
	// userInput , _ := reader.ReadString('.')
	// fmt.Println(userInput)


	fmt.Println("Enter ur email")
	// var name string

	// fmt.Scan(&name)

	// fmt.Println(name)


	// bufio 
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('.')

	fmt.Println(value)




}