package main

import (
	"fmt"
	"strconv"
)


func main(){

	// string to int conversion
	str := "95"
    
	strToInt , _ := strconv.Atoi(str)

	fmt.Printf("value is %d and type is %T\n", strToInt, strToInt)


	// int to string

	num := 45

	numToStr := strconv.Itoa(num)

	fmt.Printf("value is %s, type is %T\n", numToStr, numToStr)

	// string to float 
	flt1 := "756.465"

	strtoFlt, _ := strconv.ParseFloat(flt1, 64)
	fmt.Printf("value is %f and type is %T\n", strtoFlt, strtoFlt)
 

	// float to string 
	flt := 75.590
	fltToStr := strconv.FormatFloat(flt, 'f', 3, 64)
	fmt.Printf("Type of %s is %T\n", fltToStr, fltToStr)

}