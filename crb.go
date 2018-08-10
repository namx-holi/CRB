

package main

import (
	"fmt"
	// "net/http"
	// "time"
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

	fmt.Printf("URL to request: %s\n", *urlPtr)
	fmt.Printf("We will hit this with %d hits\n", *countPtr)
	fmt.Printf("We will repeat this process %d times\n", *loopPtr)
	fmt.Printf("We will wait %dms between the processes\n", *cooldownPtr)

	if *verbosePtr {
		fmt.Println("We are running in verbose mode")
	} else {
		fmt.Println("We are not running in verbose mode")
	}
}
