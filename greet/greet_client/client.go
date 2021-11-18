package main

import (
	"context"
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

	unary(client)

}

func unary(client greetpb.GreetServiceClient) {

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
