package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Backend received request: %s %s\n", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello from Backend! Path: %s", r.URL.Path)
	})

	fmt.Println("Dummy Backend listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
