package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bensooraj/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial/connect: %v\n", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created the client: %f", c)

	doUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Unary RPC initiated")

	req := &calculatorpb.SumRequest{
		FirstNumber:  123,
		SecondNumber: 321,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v \n", err)
	}

	log.Printf("Response: %v\n", res.SumResult)
}
