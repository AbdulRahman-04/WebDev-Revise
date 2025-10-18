package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=========================
CHANNEL BASICS
=========================
- Channel = pipe, used to communicate between goroutines
- Sender writes, receiver reads
- Writing blocks the sender until another goroutine reads (unbuffered)
- Reading blocks the receiver until sender writes
- Deadlock occurs if no goroutine is ready for opposite operation
- main() is the entry point, cannot be called with arguments or as goroutine
*/

// ==========================
// 1️⃣ Simple Unbuffered Channel
// ==========================
func SimpleUnbuffered() {
	ch := make(chan string) // unbuffered

	go func() {
		msg := <-ch
		fmt.Println("Received:", msg)
	}()

	ch <- "Hello from unbuffered channel"
	// Sender blocks until receiver reads
	time.Sleep(100 * time.Millisecond) // just to ensure print
}

// ==========================
// 2️⃣ Buffered Channel Example
// ==========================
func BufferedExample() {
	ch := make(chan string, 3) // buffer size 3

	ch <- "Hello"
	ch <- "World"
	ch <- "Go" // all sent without blocking

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	/*
		💡 Buffered channel analogy:
		- Channel = pipe
		- Buffer = basket inside pipe
		- Multiple items can be stored in buffer before receiver reads
		- Sender blocks only when buffer is full
	*/
}

// ==========================
// 3️⃣ Channel with Range and Close
// ==========================
func ChannelRangeClose() {
	ch := make(chan string)

	go func() {
		for val := range ch {
			fmt.Println("Range received:", val)
		}
		fmt.Println("Channel closed ✅")
	}()

	ch <- "hello from channel"
	ch <- "kysa h bro?"
	ch <- "kidr jana h teku?"
	ch <- "chal chod detao teku"

	close(ch) // closing channel stops the range loop
	time.Sleep(100 * time.Millisecond)
}

// ==========================
// 4️⃣ Directional Channels (Send-only / Receive-only)
// ==========================
func DirectionalChannels() {
	ch := make(chan int)

	// Receiver goroutine
	go func(rcv <-chan int) {
		fmt.Println("Received:", <-rcv)
	}(ch)

	// Sender goroutine
	go func(snd chan<- int) {
		snd <- 7878
	}(ch)

	ch <- 5049 // main sends value

	time.Sleep(500 * time.Millisecond)
}

// ==========================
// 5️⃣ Buffered Channel with WaitGroup
// ==========================
func BufferedWithWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan string, 3) // buffered channel

	go func(ch chan string, wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println(<-ch)
		fmt.Println(<-ch)
		fmt.Println(<-ch)
	}(ch, &wg)

	ch <- "Hello"
	ch <- "world"
	ch <- "go"

	wg.Wait() // wait for goroutine to finish
}

// ==========================
// Main Function
// ==========================
func main() {
	fmt.Println("=== Simple Unbuffered Channel ===")
	SimpleUnbuffered()

	fmt.Println("\n=== Buffered Channel Example ===")
	BufferedExample()

	fmt.Println("\n=== Channel with Range and Close ===")
	ChannelRangeClose()

	fmt.Println("\n=== Directional Channels ===")
	DirectionalChannels()

	fmt.Println("\n=== Buffered Channel with WaitGroup ===")
	BufferedWithWaitGroup()
}

/*
=========================
KEY POINTS SUMMARY
=========================
1. Unbuffered Channel:
   - Buffer size = 0
   - Synchronous (sender & receiver must be ready)
   - One value at a time
   - Writing blocks until read

2. Buffered Channel:
   - Buffer size > 0
   - Asynchronous up to buffer limit
   - Multiple values can be stored
   - Sender blocks only if buffer full
   - Receiver blocks only if buffer empty

3. Closing a channel:
   - Allows range loops to end
   - Cannot send on closed channel

4. Directional channels:
   - chan<- T = send-only
   - <-chan T = receive-only

5. WaitGroup:
   - Ensures goroutine finishes before main exits
*/
