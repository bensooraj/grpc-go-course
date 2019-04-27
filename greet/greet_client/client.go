package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)

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
