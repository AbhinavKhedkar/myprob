package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	http.ListenAndServe(":8000", nil)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleSortRequest(w, r, false)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleSortRequest(w, r, true)
}

func handleSortRequest(w http.ResponseWriter, r *http.Request, concurrent bool) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var sortedArrays [][]int
	if concurrent {
		sortedArrays = sortConcurrently(req.ToSort)
	} else {
		sortedArrays = sortSequentially(req.ToSort)
	}

	timeTaken := time.Since(startTime).Nanoseconds()

	resp := SortResponse{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func sortSequentially(arrays [][]int) [][]int {
	var result [][]int
	for _, arr := range arrays {
		sortedArr := make([]int, len(arr))
		copy(sortedArr, arr)
		sort.Ints(sortedArr)
		result = append(result, sortedArr)
	}
	return result
}

func sortConcurrently(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	var result [][]int
	var mu sync.Mutex

	for _, arr := range arrays {
		wg.Add(1)
		go func(arr []int) {
			defer wg.Done()

			sortedArr := make([]int, len(arr))
			copy(sortedArr, arr)
			sort.Ints(sortedArr)

			mu.Lock()
			result = append(result, sortedArr)
			mu.Unlock()
		}(arr)
	}

	wg.Wait()
	return result
}
