package main

import (
	"fmt"
	"log"
	"net"

	"github.com/bensooraj/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {
	fmt.Println("Blog Server Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
