package main

import "fmt"

func main(){
 
	// creating an array 

	var arr = [3]bool {true, false , true}
	fmt.Println(arr)

	var myArr = [...]string {"hey", "hi", "sneha"}
	fmt.Println(myArr)

	fmt.Println(arr[0], arr[1])
	fmt.Println(myArr[0], myArr[2])

	fmt.Printf("length of my array is : %d\n", len(myArr))


}