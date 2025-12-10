package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Record struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
}

func generateRecords(count int) []Record {
	records := make([]Record, count)
	for i := 0; i < count; i++ {
		records[i] = Record{
			ID:        i + 1,
			Timestamp: time.Now(),
			Data: fmt.Sprintf("Data for record %d: This is some filler data to increase payload size.", i+1),
		}
	}
	return records
}

func handleBackend1(w http.ResponseWriter, r *http.Request) {
	const recordCount = 10000

	startTime := time.Now()
	records := generateRecords(recordCount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":          "success",
		"record_count":    recordCount,
		"records":         records,
		"generation_time": time.Since(startTime).String(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing JSON response: %v", err)
		http.Error(w, "Error creating response", http.StatusInternalServerError)
	}

	log.Printf("Successfully served %d records in %s", recordCount, time.Since(startTime))
}

func main() {
	http.HandleFunc("/data-service", handleBackend1)

	port := ":8081"
	log.Printf("Starting high-payload backend on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}