package main

import (
	"errors"
	"fmt"
)

func OddEven(x int, y int) (int, error){
  if x %2 != 0 || y%2!= 0{
    return  0, errors.New("x or y is an odd number")
  }
 
  return  x + y, nil

}


func main(){
 
  result, err := OddEven(12, 15)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(result)
  }


}