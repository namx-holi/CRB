

package main

import (
	"fmt"
	"net/http"
	"time"
	// "math"
	"flag"
)


func main() {

	// flags
	urlPtr := flag.String("url", "", "URL to request")
	countPtr := flag.Int("count", 1, "number of desired concurrent hits")
	loopPtr := flag.Int("loops", 1, "number of times to repeat the process")
	cooldownPtr := flag.Int("cooldown", 1000, "ms to wait between loops")
	verbosePtr := flag.Bool("verbose", false, "enable extra information")
	flag.Parse()

	// if no URL specified
	if (*urlPtr == "") {
		fmt.Println("Please enter a URL to hit with: -url \"https://...\"")
		return
	}

	fmt.Printf("URL to request: %s\n", *urlPtr)
	fmt.Printf("We will hit this with %d hits\n", *countPtr)
	fmt.Printf("We will repeat this process %d times\n", *loopPtr)
	fmt.Printf("We will wait %dms between the processes\n", *cooldownPtr)

	if *verbosePtr {
		fmt.Println("We are running in verbose mode")
	} else {
		fmt.Println("We are not running in verbose mode")
	}

	// channel for request times
	// size is the number of requests + 1 as to not be blocked
	request_times := make(chan float64, *countPtr + 1)

	// build request
	req, err := http.NewRequest("GET", *urlPtr, nil)
	if err != nil {
		fmt.Printf("NewRequest: %s", err)
		return
	}
}


func time_response(req_times chan<- float64, req* http.Request) {
	// http client
	client := &http.Client{}

	// make request and time it
	start := time.Now()
	_, err := client.Do(req)
	elapsed := time.Since(start)

	// if there was a request error
	if err != nil {
		fmt.Printf("Do: %s\n", err)
	}

	req_times <- elapsed.Seconds() * 1000
}