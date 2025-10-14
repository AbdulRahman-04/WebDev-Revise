package main

import "fmt"

func main(){
 
	//  classic for loop
	for i := 0; i < 12; i++ {
		// fmt.Println(i)
	}

	// while style for loop
	j := 2;
	for j <= 5 {
		// fmt.Println(j)
		j++
	}

	// break
	for  x := 2 ; x<=20; x++{
		if x <= 5 {
			break
		}
		// fmt.Println(x)
	}

	for y:= 3; y<=7; y++ {
		if y == 4 {
			continue
		}
		// fmt.Println(y)
	}


	// Range keyword in go is used for looping over slices and array elements
	// range returns two things index, value from slice or array

	myNums := make([]int, 0, 10)
	values := append(myNums, 10, 20, 30, 40, 50)

	for index, value := range values{
		fmt.Printf("index is %d and value is %d\n", index, value)
	}


}