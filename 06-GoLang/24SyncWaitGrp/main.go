package main

import (
	"fmt"
	"sync"

)
// 4. recive kro mem address ku ek wg name variable m then *syncWt mtlb btao k ye var mem address ek syncWG type ka h
func printHi(a string, wg *sync.WaitGroup){
  // 5. close the wg
  defer wg.Done() 
  fmt.Println("mEssage 1:", a)
}
func printHello(b string, wg *sync.WaitGroup){
	defer wg.Done()
   fmt.Println("Message 2:", b)
}

func main(){
	// 1 Create variable type syncWaitgroup
	var wg sync.WaitGroup

	// for loop for hi 
	for i:= 0; i<=3; i++{
		// 2. go routine jaha b h unse pehle synwt var .add krke go ku bolo k iske niche ek go routine start hona wala h, tum track kro usku
		wg.Add(1)
		// 3. syncwt type var ka mem add pass kro upar func ku taaki une uske og value ku leke methods unlock kr ske dusre function m not in main
		go printHi("hi", &wg)
	}

	// for loop for hello
	for j:= 0; j<=3; j++{
		// 2. go routine jaha b h unse pehle synwt var .add krke go ku bolo k iske niche ek go routine start hona wala h, tum track kro usku
		wg.Add(1)
		go printHello("Hello", &wg)
	}

	wg.Wait()

	
}
