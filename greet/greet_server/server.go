package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/saravase/golang_microservice_master/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	log.Printf("Greet function was invoked: %v\n", req)

	first_name := req.Greeting.FirstName
	last_name := req.Greeting.LastName

	res := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hi %s %s", first_name, last_name),
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	log.Printf("GreetManyTimes function was invoked: %v\n", req)

	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hi %s %s with Number : %d", first_name, last_name, i),
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	log.Printf("LongGreet function was invoked")

	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while streaming data from client: %v\n", err)
		}

		result = result + fmt.Sprintf("Name : %s %s\n",
			req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName())
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {

	log.Printf("GreetEveryone function was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error while reading data from client: %v\n", err)
			return err
		}

		first_name := req.GetGreeting().GetFirstName()
		last_name := req.GetGreeting().GetLastName()

		res := &greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hi %s %s", first_name, last_name),
		}

		err = stream.Send(res)
		if err != nil {
			log.Printf("Error while sending data to client: %v\n", err)
			return err
		}
	}
}

func (s *server) GreetWithTimeout(ctx context.Context, req *greetpb.GreetWithTimeoutRequest) (*greetpb.GreetWithTimeoutResponse, error) {

	log.Printf("GreetWithTimeout function was invoked: %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "Request reached timeout")
		}
		time.Sleep(1 * time.Second)
	}

	first_name := req.Greeting.FirstName
	last_name := req.Greeting.LastName

	res := &greetpb.GreetWithTimeoutResponse{
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
