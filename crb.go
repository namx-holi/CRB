

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
	// verbosePtr := flag.Bool("verbose", false, "enable extra information")
	flag.Parse()

	// if no URL specified
	if (*urlPtr == "") {
		fmt.Println("Please enter a URL to hit with: -url \"https://...\"")
		return
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
	for loopNb := 1; loopNb <= *loopPtr; loopNb++ {
		if *loopPtr > 1 {
			fmt.Printf("--LOOP %d OF %d--\n", loopNb, *loopPtr)
		}

		all_times[loopNb-1] = benchmark_process(req, *countPtr)

		time.Sleep(time.Duration(*cooldownPtr) * time.Millisecond)
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
		req_times[j] = <- req_time_chan
	}

	close(req_time_chan)

	return req_times
}
