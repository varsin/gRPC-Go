package main

import (
	"grpc-go-framework/video/videopb"
	"io"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct{}

var count int

func (*server) VideoStream(stream videopb.VideoService_VideoStreamServer) error {
	log.Printf("Bidirectional streaming has been started!!")
	for {
		count++
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Client has stopped streaming any data")
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading client data %v", err)
			return err
		}
		name := req.GetVideo().GetName()
		result := "Got Frame " + strconv.Itoa(count) + " " + name
		err = stream.Send(&videopb.VideoStreamResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while streaming data from server %v", err)
			return err
		}
	}
}

func main() {
	// Here we are doing Port Binding
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Creating grpc server.
	s := grpc.NewServer()
	videopb.RegisterVideoServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}
}
