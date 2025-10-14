package main

import "fmt"

type Student struct {
	name string
	age int
	Area string
	isPass bool
}


func main(){

	studDetails := Student {
		name: "Rahman",
		age: 21,
		Area: "Chandulal Baradari",
		isPass: true,
	}

	fmt.Println(studDetails)
	fmt.Println(studDetails.name, studDetails.age)

}