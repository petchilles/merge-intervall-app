package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// a single interval
type Interval struct {
	Start int
	End   int
}

func main() {
	// Parse command-line arguments
	numIntervals := flag.Int("n", 1000, "number of intervals to generate")
	flag.Parse()

	// Create a new random generator with a seed based on the current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	intervals := make([]Interval, *numIntervals)

	for i := 0; i < *numIntervals; i++ {
		start := rng.Intn(*numIntervals * 20) - *numIntervals * 10 // random number between -10000 and 10000
		end := start + rng.Intn(100)
		intervals[i] = Interval{
			Start: start,
			End:   end,
		}
	}

	file, err := os.Create("test_intervals.txt")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	for _, interval := range intervals {
		// Write the interval to the file in the specified format
		if _, err := fmt.Fprintf(file, "[%d,%d]", interval.Start, interval.End); err != nil {
			fmt.Println("Failed to write interval:", err)
			return
		}
	}

	fmt.Printf("Generated %d intervals and saved to test_intervals.txt\n", *numIntervals)
}
