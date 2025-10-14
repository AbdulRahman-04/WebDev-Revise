package main

import "fmt"

func main(){
 
	temp := 12

	switch temp {
	case 1,2,3,4,5,6,7,8,9,10 :
		fmt.Println("Very cold")
	case 11, 12, 13, 14, 15 ,16 , 17, 18, 19, 20:
		fmt.Println("cold")
	case 21, 22, 23, 24, 25, 26, 27, 28, 29, 30:
		fmt.Println("Really great weather")
	default: 
	fmt.Println("extremely hot")			
	}


	day := 2

	switch day {
	case 1: 
	fmt.Println("sunday")
	case 2:
		fmt.Println("monday")
	case 3: 
	fmt.Println("tuesday")
	case 4: 
	fmt.Println("wednesday")
	case 5:
		fmt.Println("thursday")
	case 6:
		fmt.Println("friday")
	case 7:
		fmt.Println("saturday")			
	}

}