package main

import (
	"fmt"
	"strings"
)

func main() {

	// 1️⃣ strings.Split – Break a string into pieces

	str1 := "hello babe kaise ho aap"

	words := strings.Split(str1, " ")
	fmt.Println(words)

	// 2️⃣ strings.Count – Count how many times a substring appears
	fruit := "banana"

	myCount := strings.Count(fruit, "a")
	fmt.Println(myCount)


   // 3️⃣ strings.TrimSpace – Remove leading/trailing spaces
   myVar := "        hey bro"
   myTrim := strings.TrimSpace(myVar)
   fmt.Println(myTrim) 


  
   // 4️⃣ strings.Join – Join a slice of words 
   parts := []string {"hi", "my", "bro"}

   joinIt := strings.Join(parts, " ")
   fmt.Println(joinIt)

   
//     // 5️⃣ strings.HasPrefix – Check if string starts with prefix
    file := "log_2025_error.txt"
    hasLog := strings.HasPrefix(file, "log")
    fmt.Println("5. Starts with 'log'? :", hasLog)

}
