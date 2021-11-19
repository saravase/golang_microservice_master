package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/saravase/golang_microservice_master/player/playerpb"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while connecting grpc server: %v\n", err)
	}

	client := playerpb.NewPlayerServiceClient(conn)

	doVideoStreaming(client)

}

func doVideoStreaming(client playerpb.PlayerServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 31*time.Second)
	defer cancel()

	stream, err := client.GetVideo(ctx, &playerpb.Request{Filename: "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4"})
	if err != nil {
		log.Fatalf("Error while calling GetVideo server streaming gRPC: %v\n", err)
	}

	file, err := os.Create("demo.mp4")
	if err != nil {
		log.Fatalf("Error while creating file: %v\n", err)
	}

	defer file.Close()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming data from server: %v\n", err)
		}

		file.Write(res.GetContent())
	}
}
