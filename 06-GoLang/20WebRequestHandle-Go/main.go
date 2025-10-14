package main

import (
	"fmt"
	"io"
	"net/http"
)

func main(){
	url := "https://jsonplaceholder.typicode.com/posts/1"

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}