package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/saravase/golang_microservice_master/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doClientStreaming(client)

	// doBiDirectionalStreaming(client)

	doUnaryWithTimeout(client, 5*time.Second)

	doUnaryWithTimeout(client, 1*time.Second)

}

func doBiDirectionalStreaming(client greetpb.GreetServiceClient) {

	log.Println("GreetEveryone bi-directional streaming RPC triggered ....")

	greetings := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "saravana",
				LastName:  "kumar",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "optimus",
				LastName:  "prime",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "prabhu",
				LastName:  "deva",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "chithra",
				LastName:  "deva",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "kasthoori",
				LastName:  "maan",
			},
		},
	}

	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone bi-directional streaming RPC: %v\n", err)
	}

	ch := make(chan struct{})

	go func() {
		for _, greet := range greetings {
			err := stream.Send(greet)
			if err != nil {
				log.Fatalf("Error while sending to server: %v\n", err)
			}
			log.Printf("Sending Greet : %v\n", greet)
			time.Sleep(1111 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(ch)
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming data from server: %v\n", err)
				close(ch)
			}

			log.Printf("Server response : %v\n", res.GetResult())
		}
	}()

	<-ch
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

func doUnaryWithTimeout(client greetpb.GreetServiceClient, timeout time.Duration) {

	log.Println("GreetWithTimeout unary RPC triggered ....")

	req := &greetpb.GreetWithTimeoutRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Optimus",
			LastName:  "Primz",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.GreetWithTimeout(ctx, req)
	if err != nil {
		status, ok := status.FromError(err)

		if ok {
			if status.Code() == codes.DeadlineExceeded {
				log.Fatalln("Timeout exceeded")
			} else {
				log.Fatalf("Other error occurred: %v\n", err)
			}
		} else {
			log.Fatalf("Couldn't able to call greet unary RPC: %v\n", err)
		}

		return
	}

	log.Printf("Result greet unary RPC: %v\n", res)
}
