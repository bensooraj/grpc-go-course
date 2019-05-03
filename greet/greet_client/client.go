package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bensooraj/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial/connect: %v ", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created the client: %f", c)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiDiStreaming(c)

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
