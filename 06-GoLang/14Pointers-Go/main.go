package main

import "fmt"


func main(){
	myNum := 49

	myPtr := &myNum
   
	fmt.Println( *myPtr)

	modifyValueByReference(myPtr)

	fmt.Println(*myPtr)

}


func modifyValueByReference(myPtr *int){
	*myPtr = 5049
}