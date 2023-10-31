package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	addr = flag.String("addr", "http://localhost:8080/comments", "the address of service 2 to get data")
)

func main() {
	flag.Parse()

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan []byte, 1) // Create a channel to collect responses

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", *addr, nil)
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				ch <- nil
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Failed to get data: %v", err)
				ch <- nil
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response body: %v", err)
				ch <- nil
				return
			}

			ch <- data
		}()

		select {
		case data := <-ch:
			if data == nil {
				http.Error(w, "Failed to get data", http.StatusInternalServerError)
				return
			}
			// fmt.Println(string(data))

			var comments []map[string]interface{}
			if err := json.Unmarshal(data, &comments); err != nil {
				http.Error(w, "Failed to unmarshal data", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("Received all comments. Check server logs for details."))
		case <-time.After(5 * time.Second):
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
		}
	})

	log.Println("HTTP server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
