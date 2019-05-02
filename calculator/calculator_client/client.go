package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client Stream RPC initiated")
	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 5,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 6,
		},
	}

	clientStream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage RPC: %v \n", err)
	}

	for _, req := range requests {
		clientStream.Send(req)
		fmt.Println("Sent number: ", req.GetNumber())
		time.Sleep(500 * time.Microsecond)
	}

	res, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v \n", err)
	}
	fmt.Println("Average: ", res.GetAverageResult())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Server Stream RPC initiated")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12881624,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v \n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("Done receiving prime factors: %v \n", err)
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving prime factors: %v \n", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
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
