package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Backend example starting on :8080")
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
