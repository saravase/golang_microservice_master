package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/saravase/golang_microservice_master/player/playerpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) GetVideo(req *playerpb.Request, stream playerpb.PlayerService_GetVideoServer) error {

	res, err := http.Get(req.GetFilename())
	if err != nil {
		log.Fatalf("Error while requesting file: %v\n", err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error while reading file: %v\n", err)
	}

	err = stream.Send(&playerpb.Response{
		Content: data,
	})
	if err != nil {
		log.Fatalf("Error while sending date to client: %v\n", err)
	}

	return nil

}

func main() {

	listener, err := net.Listen("tcp", "localhost:50053")
	if err != nil {
		log.Fatalf("Error while creating listener: %v\n", err)
	}

	s := grpc.NewServer()

	playerpb.RegisterPlayerServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error while serving : %v\n", err)
	}

	log.Printf("Listening on port: %v\n", listener.Addr().String())

}
