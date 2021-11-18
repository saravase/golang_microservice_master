package main

import (
	"context"
	"io"
	"log"

	"github.com/saravase/golang_microservice_master/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	// Dial TCP
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect tcp: %v\n", err)
	}

	defer conn.Close()

	// Register Client
	client := greetpb.NewGreetServiceClient(conn)

	log.Printf("gRPC greet client: %v\n", client)

	// doUnary(client)

	doServerStreaming(client)

}

func doServerStreaming(client greetpb.GreetServiceClient) {

	log.Println("GreetManyTimes server streaming RPC triggered ....")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Optimus",
			LastName:  "Primz",
		},
	}

	stream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't able to call greet unary RPC: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}

		log.Printf("Response from GreetManyTimes RPC: %v\n", res.GetResult())
	}
}

func doUnary(client greetpb.GreetServiceClient) {

	log.Println("Greet unary RPC triggered ....")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Optimus",
			LastName:  "Primz",
		},
	}

	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't able to call greet unary RPC: %v\n", err)
	}

	log.Printf("Result greet unary RPC: %v\n", res)
}
