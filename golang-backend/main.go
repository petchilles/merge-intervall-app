/*
This package handles an http request carrying a slice of intervals to be processed.
Overlapping intervals are merged. The http response delivers the merged intervals.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/rs/cors"
)

// configuration parameters
type Config struct {
	MaxIntervals int `json:"max_intervals"`
}

// a single interval
type Interval struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// incoming type
type MergeRequest struct {
	Intervals []Interval `json:"intervals"`
}

// outgoing type
type MergeResponse struct {
	Result       []Interval `json:"result"`
	ElapsedTime  string     `json:"elapsed_time"`
	MemoryUsage  string     `json:"memory_usage"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}

// merges intervals
func mergeIntervals(intervals []Interval) []Interval {
	// pass through empty input array
	if len(intervals) == 0 {
		return []Interval{}
	}

	// sort intervals by start, end ascending
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].Start == intervals[j].Start {
			return intervals[i].End < intervals[j].End
		}
		return intervals[i].Start < intervals[j].Start
	})

	// initialise to be returned "merged" slice
	// with the first item from intervals
	merged := []Interval{intervals[0]}

	// loop sorted intervals starting from the second item
	for _, current := range intervals[1:] {
		last := &merged[len(merged)-1]

		// compare current item with last item of "merged"
		if current.Start <= last.End {
			if current.End > last.End {
				// we have an overlap, therefore increase the end
				// of the last item in "merged"
				last.End = current.End
			}
		} else {
			// no overlap, therefore append the current item to "merged"
			merged = append(merged, current)
		}
	}

	return merged
}

// http handler
func mergeHandler(w http.ResponseWriter, r *http.Request) {
	var req MergeRequest
	config := r.Context().Value("config").(Config)

	// decode json data from request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// reject interval counts greater than defined in the config file
	if len(req.Intervals) > config.MaxIntervals {
		http.Error(w, fmt.Sprintf("Too many intervals. Maximum allowed number of intervals: %d", config.MaxIntervals), http.StatusRequestEntityTooLarge)
		return
	}

	// record time and memory usage
	startTime := time.Now()
	var memStatsStart runtime.MemStats
	runtime.ReadMemStats(&memStatsStart)

	// process decoded data using the mergeIntervals function
	result := mergeIntervals(req.Intervals)

	// measure elapsed time and memory usage
	elapsedTime := time.Since(startTime)
	var memStatsEnd runtime.MemStats
	runtime.ReadMemStats(&memStatsEnd)
	memoryUsage := memStatsEnd.TotalAlloc - memStatsStart.TotalAlloc

	resp := MergeResponse{
		Result:      result,
		ElapsedTime: elapsedTime.String(),
		MemoryUsage: fmt.Sprintf("%d bytes", memoryUsage),
	}

	// json encode and send response data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// entry point
func main() {
	// Set configuration for maximum number of intervals
	myConfig := Config{MaxIntervals: 1000000}

	// Create a ServeMux
	mux := http.NewServeMux()

	// Create a handler that saves the configuration in the request context and calls mergeHandler
	mux.HandleFunc("/merge", func(w http.ResponseWriter, r *http.Request) {
		// Save configuration in the request context
		ctx := context.WithValue(r.Context(), "config", myConfig)
		// Create a new request with the updated context
		r = r.WithContext(ctx)
		// Call MergeHandler
		mergeHandler(w, r)
	})

	// Wrap the default mux with the CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3010", "http://127.0.0.1:3010"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := corsHandler.Handler(mux)

	log.Println("Server running on port 8085")
	if err := http.ListenAndServe(":8085", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
