package main

import "fmt"

func main(){
	// creating a slice 
	mySlice := []int {49, 46, 42}
	fmt.Println(mySlice)

	myslice2 := []bool {true, false}
	fmt.Println(myslice2)

	myslice3 := []string {"fahad", "syed", "omer"}
	fmt.Println(myslice3)

	
 // using slice we can also extract values from an array just like js

    myslice4 := []int {49, 52, 55, 69}
	slicedSlice := myslice4[0:3]
	fmt.Println(slicedSlice)

    myslice5 := []string {"hey", "hi", "faisal"}
	fmt.Println(myslice5)


	
	// // USING MAKE FUNCTION TO CREATE SLICES
	// /*
	//  "Length" bataata hai initially kitne elements stored h slice m and 
	//   "capacity" maximum limit set karta hai slice ki.
	//   Agar tu make([]int, 5, 10) karega, toh 5 elements initially hain,
	//    par slice mein maximum 10 elements tak ja sakte hain
	// */

	// newSlice := make([]bool , 0, 10)
	// values := append(newSlice, true, false , false, true, true, false)
	// fmt.Println(values)

	// newSlice := make([]string , 0 , 5)
	// values := append(newSlice, "rxhman", "dev", "bs", "fck")
	// fmt.Println(values)
    
	// newSlice := make([]int, 0, 6)
	// values := append(newSlice, 5049, 5046, 5042)
	// fmt.Println(values)

}

