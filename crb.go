
package main

import (
	"fmt"
	"net/http"
	"time"
	// "math"
	"flag"
	"strings"
	"sort"
)

const INDENTATION_WIDTH int = 2

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

	results := run_benchmark(*urlPtr, *loopPtr, *countPtr, *cooldownPtr, *verbosePtr)

	compiled_results := compile_results(results)

	fmt.Println(compiled_results)
}

func run_benchmark(url string, loops int, count int, cooldown int, verbose bool) ([][]float64) {
	// build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("NewRequest error: %s\n", err)
	}

	// init array for storing all the stats
	results := make([][]float64, loops)

	// number of processes
	for loopNb := 1; loopNb <= loops; loopNb++ {
		if loops > 1 {
			fmt.Printf("\n--LOOP %d OF %d--\n", loopNb, loops)
		}

		result := benchmark_process(req, count)

		if verbose {
			fmt.Printf("\n\n  --STATS FOR LOOP %d--\n", loopNb)
			display_stats(result, 2)

		}
		fmt.Printf("\n\n")

		results[loopNb-1] = result

		if loopNb != loops {
			time.Sleep(time.Duration(cooldown) * time.Millisecond)
		}
	}

	return results
}

func benchmark_process(req* http.Request, count int) ([]float64){
	// channel and array for request times
	// channel is used because it's quick
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

func compile_results(results [][]float64) ([]float64) {
	var compiled_results []float64

	for _, result_set := range results {
		for _, result := range result_set {
			compiled_results = append(compiled_results, result)
		}
	}

	// sort em all
	sort.Float64s(compiled_results)

	return compiled_results
}

func calculate_min(req_times []float64) (float64) {
	var min_value float64 = req_times[0]

	for i := 1; i < len(req_times); i++ {
		if req_times[i] < min_value {
			min_value = req_times[i]
		}
	}

	return min_value
}

func calculate_max(req_times []float64) (float64) {
	var max_value float64 = req_times[0]

	for i := 1; i < len(req_times); i++ {
		if req_times[i] > max_value {
			max_value = req_times[i]
		}
	}

	return max_value
}

func calculate_mean(req_times []float64) (float64) {
	var total_value float64 = 0

	for i := 0; i < len(req_times); i++ {
		total_value += req_times[i]
	}

	return total_value / float64(len(req_times))
}

func calculate_median(req_times []float64) (float64) {
	n := len(req_times)

	if (n % 2) == 0 {
		// mean of values either side of the middle
		return (req_times[n/2 - 1] + req_times[n/2])/2.0
	}
	return req_times[(n-1)/2]
}

func display_stats(req_times []float64, indent int) {
	indentation := strings.Repeat(" ", indent*INDENTATION_WIDTH)
	fmt.Printf("%sMin response time was:    %fms\n", indentation, calculate_min(req_times))
	fmt.Printf("%sMean response time was:   %fms\n", indentation, calculate_mean(req_times))
	fmt.Printf("%sMax response time was:    %fms\n", indentation, calculate_max(req_times))
	fmt.Printf("%sMedian response time was: %fms\n", indentation, calculate_median(req_times))
}
