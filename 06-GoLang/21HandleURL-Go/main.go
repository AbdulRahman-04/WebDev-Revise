package main

import (
	"fmt"
	"net/url"
)

func main(){
	myUrl := "https://github.com/AbdulRahman-04?tab=repositories"

	urlConv , err := url.Parse(myUrl)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Type is %T\n", urlConv)
	fmt.Println("scheme:", urlConv.Scheme)
	fmt.Println("host", urlConv.Host)
	fmt.Println("path", urlConv.Path)
	fmt.Println("query", urlConv.Query())
}