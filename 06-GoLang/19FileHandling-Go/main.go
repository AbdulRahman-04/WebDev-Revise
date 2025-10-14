package main

import (
	"fmt"
	"os"
)

func main(){

	// creating a file
	fileCreate, err := os.Create("sendEmail.txt")
	if err != nil {
		fmt.Println("err creating a file")
	}
	fmt.Println("file created with name :", fileCreate.Name())

	// writing content into file
	_, err = fileCreate.WriteString("Hello bhai\n kidr h ki?") 
    if err != nil {
		fmt.Println("error writing content in file")
	}
	fmt.Println("content added to file✅")

	// open file for reading
	openFile, err := os.Open(fileCreate.Name())
	if err != nil {
		fmt.Println("error opening file")
	}

	defer openFile.Close()

		// 🔹 Step 4: Create a buffer to hold file content
	buffer := make([]byte, 1024) // 1KB buffer

	// 🔹 Step 5: Read content into buffer
	bytesRead, err := openFile.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading file")
		return
	}

	// 🔹 Step 6: Print the content read
	fmt.Println("Bytes read:", bytesRead)
	fmt.Println("File content:\n" + string(buffer[:bytesRead]))
}


