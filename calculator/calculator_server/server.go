package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/bensooraj/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
}

func (*server) ComputeAverage(clientStream calculatorpb.CalculatorService_ComputeAverageServer) error {

	sum := int64(0)
	count := int64(0)

	for {
		req, err := clientStream.Recv()
		if err == io.EOF {
			log.Printf("Reached end of stream EOF: %v ", err)
			return clientStream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				AverageResult: float32(sum) / float32(count),
			})
		}
		if err != nil {
			log.Fatalf("Error receiving from the client stream: %v ", err)
		}

		num := req.GetNumber()
		sum += num
		count++
	}
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	fmt.Printf("Received Sum RPC: %v\n", req)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	fmt.Println("Server Stream PrimeNumberDecomposition initiated")

	primeNumber := req.GetNumber()
	divisor := int64(2)

	for primeNumber > 1 {
		if primeNumber%divisor == 0 {

			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			primeNumber = primeNumber / divisor

		} else {
			divisor = divisor + 1
			fmt.Println("Divisor incremented to : ", divisor)
		}
	}

	return nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
