package main

import (
	"context"
	"io"
	"log"

	"github.com/saravase/golang_microservice_master/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect : %v\n", err)
	}

	client := calculatorpb.NewCalculatorServiceClient(conn)

	// doUnary(client)

	// doServerStreaming(client)

	doClientStreaming(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.CalculatorRequest{
		Numbers: &calculatorpb.Numbers{
			Num1: 3,
			Num2: 10,
		},
	}

	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't call unary sum RPC : %v\n", err)
	}

	log.Printf("Sum of %d and %d : %d\n", req.Numbers.Num1, req.Numbers.Num2, res.Result)
}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	stream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't call PrimeNumberDecomposition server streaming RPC : %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while streaming data from server: %v\n", err)
		}
		log.Printf("Primes: %v\n", res.GetPrime())
	}
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {

	numArr := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}

	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Printf("Error while calling ComputeAverage client streaming RPC: %v\n", err)
	}

	for _, val := range numArr {
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: val,
		})

		if err != nil {
			log.Printf("Error while sending data to server: %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error while receving data from server: %v\n", err)
	}

	log.Printf("Average : %v\n", res.Average)
}
