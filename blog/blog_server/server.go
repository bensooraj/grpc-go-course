package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/bensooraj/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

type server struct {
}

func main() {

	// In case of crash, get the filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Server Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// Wait for Ctrl-C to exit
	osSignalChannel := make(chan os.Signal, 1)
	signal.Notify(osSignalChannel, os.Interrupt)

	// Block until a signal is received
	<-osSignalChannel

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End Of Program")

}
