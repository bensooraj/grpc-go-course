package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/status"

	"github.com/bensooraj/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am a client")

	certFile := "ssl/ca.crt"
	creds, sslError := credentials.NewClientTLSFromFile(certFile, "")
	if sslError != nil {
		log.Fatalf("Failed to load SSL certificates: %v", sslError)
		return
	}
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial/connect: %v ", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created the client: %f", c)

	doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	// doUnaryWithDeadline(c, 5*time.Second)
	// doUnaryWithDeadline(c, 1*time.Second)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("BiDi Streaming initiated")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ben",
				LastName:  "Sooraj",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hannah",
				LastName:  "Angeline",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saasha Mehr",
				LastName:  "Sooraj",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Surya",
				LastName:  "Mohan",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Eunice",
				LastName:  "Keren",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v \n", err)
		return
	}

	waitChannel := make(chan struct{})

	// Send messages to client
	go func() {
		for _, req := range requests {
			fmt.Println("[SENDING] Message: ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// Receive messages from client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Reached EOF")
				break
			}
			if err != nil {
				log.Fatalf("Error while reading server stream: %v ", err)
				break
			}
			fmt.Println("[RECEIVING] Message: ", res)
		}
		close(waitChannel)
	}()

	<-waitChannel
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client Streaming initiated")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ben",
				LastName:  "Sooraj",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hannah",
				LastName:  "Angeline",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saasha Mehr",
				LastName:  "Sooraj",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Surya",
				LastName:  "Mohan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Eunice",
				LastName:  "Keren",
			},
		},
	}

	clientStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v \n", err)
	}

	for _, req := range requests {
		clientStream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v \n", err)
	}
	fmt.Println("Long greet response: ", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Server Streaming initiated")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ben",
			LastName:  "Sooraj",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes rpc: %v \n", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			fmt.Println("Reached Stream EOF")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving from the server stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %v\n", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Unary RPC initiated")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ben",
			LastName:  "Sooraj",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet rpc: %v \n", err)
	}

	log.Printf("Response: %v\n", res)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Println("Unary RPC initiated")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ben",
			LastName:  "Sooraj",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusError, ok := status.FromError(err)
		if ok {
			if statusError.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline exceeded!", statusError.Message())
			} else {
				fmt.Println("Unexpected Error!", statusError)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v \n", err)
		}

	}

	log.Printf("Response: %v\n", res)
}
