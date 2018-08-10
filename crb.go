

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


	// build request
	req, err := http.NewRequest("GET", *urlPtr, nil)
	if err != nil {
		fmt.Printf("NewRequest: %s", err)
		return
	}

	// init array for storing all the stats
	all_times := make([][]float64, *loopPtr)

	// number of processes
	for loopNb := 0; loopNb < *loopPtr; loopNb++ {
		all_times[loopNb] = benchmark_process(req, *countPtr)
	}
}


func time_response(req* http.Request, req_time_chan chan<- float64) {
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

	req_time_chan <- elapsed.Seconds() * 1000
}


func benchmark_process(req* http.Request, count int) ([]float64){
	// channel for request times
	req_time_chan := make(chan float64, count)

	req_times := make([]float64, count)

	// start all requests
	for i := 0; i < count; i++ {
		go time_response(req, req_time_chan)
	}

	// collect responses
	for j := 0; j < count; j++ {
		req_times[count] = <- req_time_chan
	}

	return req_times
}
