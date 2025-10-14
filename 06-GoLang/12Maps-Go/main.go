package main

import "fmt"

func main(){

	// what is maps ?
	// maps basically ek unordered collection of data h jo key value pair m
	// data store krta and evry key unique rehti jisse uski value ko aap retrive kr skte
	// also map use krke aap data store krskte retrive krskte nd delete b kr skte

	// make function to make a map 
	myMap := make(map[string] string)

	myMap["name"] = "Rahman"
	myMap["age"] = "21"
	myMap["isAlive"] = "true"

	fmt.Println(myMap)

	// LITERAL DECLARATION(object type)
	// colors := map[string] string{
	// 	"color1": "blue",
	// 	"color2": "yellow",
	// }

	// fmt.Println(colors["color1"])
}