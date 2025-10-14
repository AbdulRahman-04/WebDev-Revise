package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Age int `json:"int"`
}

func main(){

	// MARSHALLING STRUCT TO JSON
	myDetails := Student{
		Name: "rxhman",
		Email: "rxhman87@gmail.com",
		Age: 21,
	}

	jsonData , err := json.Marshal(myDetails)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(jsonData))
 


	// UNMARSHALLING JSON TO STRUCT
	var myUnMarshal Student

	err = json.Unmarshal(jsonData, &myUnMarshal)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(myUnMarshal)
	
}