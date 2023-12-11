// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// RequestPayload represents the input JSON structure
type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

// ResponsePayload represents the output JSON structure
type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       string  `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	fmt.Println("Server listening on :8000")
	http.ListenAndServe(":8000", nil)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, sequentialSort)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, concurrentSort)
}

func handleRequest(w http.ResponseWriter, r *http.Request, sorter func([][]int) [][]int) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sorter(requestPayload.ToSort)
	endTime := time.Now()

	responsePayload := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       fmt.Sprint(endTime.Sub(startTime).Nanoseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

func sequentialSort(toSort [][]int) [][]int {
	for i := range toSort {
		sort.Ints(toSort[i])
	}
	return toSort
}

func concurrentSort(toSort [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(toSort))

	for i := range toSort {
		go func(i int) {
			defer wg.Done()
			// Create a copy of the sub-array before sorting
			temp := make([]int, len(toSort[i]))
			copy(temp, toSort[i])
			sort.Ints(temp)
			toSort[i] = temp
		}(i)
	}

	wg.Wait()
	return toSort
}
