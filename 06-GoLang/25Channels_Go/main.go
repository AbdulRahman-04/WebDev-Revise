package main

import "fmt"


/* THERE ARE TWO TYPES OF CHANNELS IN GO : BUFFERED , UNBUFFERED
   
   UNBUFFERD : YE AISA CHANNEL H JISKU CREATE KRTE WAQT AAP YE CHANNEL M KITTA DATA STORE KRSKTE WO LIMIT NI DETE. E.G
   chan1 := make(chan string) -> idhr aap sirf channel bnaadiye bina bataye k isme kitta data store hunga

   BUFFERED: YE AISA CHANNEL H JISME AAP BATATE K KITTA DATA STORE KRINGE MAXIMUM ND AGAR LIMIT SE ZYADA DATA BHEJE TOH DEADLOCK ERR AJATA.
   chan2 := make(chan string, 2) -> 2 string ka data store krsktu m sirf 

   buffered channel k data print krana h toh pehle sender func m channel k andar data daalke close krdena channel close(chan2) then main func m
   for loop m range keyword chan2 pe lgake use data print krwa skte 
   for data := range chan2  {
        fmt.println(data)
     }

	unbuffered m b aap chuncks m data send krskte ek hi channel m lekin wahi channel close krke buffered jaisa loop krake print krwana pdta data 

*/


//////////////////////////////
// ✅ CASE 1: Unbuffered channel with only ONE message
//////////////////////////////

// 🔹 Sender sends only one value
func sendSingle(ch chan string) {
	ch <- "Sirf ek message bhai"
}

//////////////////////////////
// ✅ CASE 2: Buffered/Unbuffered channel with MULTIPLE messages
//////////////////////////////

// 🔹 Sender sends multiple values and closes the channel after done
func sendMultiple(ch chan string) {
	ch <- "Message 1"
	ch <- "Message 2"
	ch <- "Message 3"
	close(ch) // 🔒 Always close the channel from sender side when done sending
}

func main() {
	// ----------- CASE 1 -------------
	ch1 := make(chan string) // Unbuffered channel
	go sendSingle(ch1)

	// ✅ For single value — direct receive into a variable
	msg := <-ch1
	fmt.Println("Single Value Received:", msg)

	// ----------- CASE 2 -------------
	ch2 := make(chan string, 3) // Buffered channel (can also be unbuffered)

	go sendMultiple(ch2)

	// ✅ For multiple values — use range (requires channel to be closed)
	for val := range ch2 {
		fmt.Println("Multiple Value Received:", val)
	}
}