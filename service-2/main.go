
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	pb "golang/grpc/proto"
)

var (
	port = flag.Int("port", 8080, "The server port") // 50051
)


type commentServer struct {
	pb.UnimplementedCommentServiceServer
}

func (s *commentServer) GetAllComments(ctx context.Context, req *pb.EmptyRequest) (*pb.CommentList, error) {
	// Fetch all comments from your data source and populate a CommentList object
	// Read JSON file
	data, err := os.ReadFile("comments.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into a slice of Comment objects
	var comments []*pb.Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}

	// Create and populate a CommentList object
	commentList := &pb.CommentList{
		Comments: comments,
	}

	return commentList, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	pb.RegisterCommentServiceServer(s, &commentServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
