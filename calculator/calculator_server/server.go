package main

import (
	"context"
	"log"
	"net"

	"github.com/saravase/golang_microservice_master/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	num1 := req.Numbers.Num1
	num2 := req.Numbers.Num2

	res := &calculatorpb.CalculatorResponse{
		Result: num1 + num2,
	}

	return res, nil
}

func main() {

	listener, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("Couldn't create listener: %v\n", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Couldn't serve: %v\n", err)
	}

	log.Printf("Server listening on port : %v\n", listener.Addr().String())

}
