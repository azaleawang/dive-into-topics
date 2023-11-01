package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("comments.json")
		if err != nil {
			http.Error(w, "Failed to read comments.json", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	log.Printf("HTTP server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
