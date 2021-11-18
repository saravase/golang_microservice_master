package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/saravase/golang_microservice_master/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	first_name := req.Greeting.FirstName
	last_name := req.Greeting.LastName

	res := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hi %s %s", first_name, last_name),
	}

	return res, nil
}

func main() {

	// Initaialize TCP
	listner, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Couldn't create listener: %v\n", err)
	}

	// Initialize gRPC server
	s := grpc.NewServer()

	// Register server
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Server
	err = s.Serve(listner)
	if err != nil {
		log.Fatalf("Couldn't serve : %v\n", err)
	}

	log.Printf("Server listening on port : %v\n", listner.Addr().String())
}
