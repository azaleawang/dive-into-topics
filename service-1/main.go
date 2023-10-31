// this service listen ti client request and get data from service 2
// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	pb "golang/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

// const (
// 	defaultName = "world"
// )

var (
	addr = flag.String("addr", "localhost:50051", "the address of service 2 to get data")
	// name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Service 1 (client) did not connect to service 2: %v", err)
	}
	defer conn.Close()

	commentClient := pb.NewCommentServiceClient(conn)

	// HTTP server
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		// Fetch all comments from Service 2
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		commentList, err := commentClient.GetAllComments(ctx, &pb.EmptyRequest{})
		if err != nil {
			http.Error(w, "Failed to get data", http.StatusInternalServerError)
			return
		}

		// Print first 5 received comments
		for i, comment := range commentList.Comments {
			if i >= 5 {
				break
			}
			log.Printf("Received Comment: %+v", comment)
		}

		// Respond to the HTTP client
		w.Write([]byte("Received all comments. Check server logs for details."))
	})
	log.Println("HTTP server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
