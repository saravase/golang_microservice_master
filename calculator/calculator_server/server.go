package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"

	"github.com/saravase/golang_microservice_master/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	var k int32 = 2
	n := req.GetNumber()

	for k > 1 {
		if n%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Prime: k,
			}
			stream.Send(res)
			n = n / k
		} else {
			k = k + 1
		}
	}

	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	var sum, c int32 = 0, 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float32(sum / c),
			})
		}

		if err != nil {
			log.Fatalf("Error while streaming data from client: %v\n", err)
		}

		c = c + 1
		sum = sum + req.GetNumber()
		log.Printf("%d %d %d", c, sum, req.GetNumber())
	}
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {

	var max int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while receving data from client: %v\n", err)
			return err
		}

		if max < req.GetNumber() {
			max = req.GetNumber()
		}
		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
			return err
		}
	}
}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	n := req.GetNumber()

	log.Printf("Number: %v\n", n)

	if n == 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Number shouldn't be 0")
	}

	return &calculatorpb.SquareRootResponse{
		Square: math.Sqrt(float64(n)),
	}, nil

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
