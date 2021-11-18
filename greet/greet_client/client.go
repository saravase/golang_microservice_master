package main

import (
	"context"
	"fmt"
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

	// doServerStreaming(client)

	doClientStreaming(client)

}

func doClientStreaming(client greetpb.GreetServiceClient) {

	log.Println("LongGreet client streaming RPC triggered ....")

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Counldn't call LongGreet client streaming RPC : %v\n", err)
	}

	for i := 0; i < 10; i++ {

		err := stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: fmt.Sprintf("%s %d", "Optimus", i),
				LastName:  fmt.Sprintf("%s %d", "Prime", i),
			},
		})

		if err != nil {
			log.Fatalf("Error, While sending data to server : %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error, While receving data from server : %v\n", err)
	}

	log.Printf("Result: %v\n", res.GetResult())
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
