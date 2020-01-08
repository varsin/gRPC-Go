package main

import (
	"context"
	"fmt"
	"grpc-go-framework/video/videopb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to the server %v", err)
	}

	defer cc.Close()

	c := videopb.NewVideoServiceClient(cc)
	bidirectionalStreaming(c)
}

func bidirectionalStreaming(c videopb.VideoServiceClient) {
	fmt.Println("Bidirectional Streaming has been started-Client")
	requests := []*videopb.VideoStreamRequest{
		&videopb.VideoStreamRequest{
			Video: &videopb.Video{
				Name: "Matrix",
			},
		},
		&videopb.VideoStreamRequest{
			Video: &videopb.Video{
				Name: "Matrix",
			},
		},

		&videopb.VideoStreamRequest{
			Video: &videopb.Video{
				Name: "Matrix",
			},
		},

		&videopb.VideoStreamRequest{
			Video: &videopb.Video{
				Name: "Matrix",
			},
		},
	}

	stream, err := c.VideoStream(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream from client side")
		return
	}

	unblock := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message from client %v\n", req)
			stream.Send(req)
			time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Streaming has been closed from server")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving content from server %v", err)
				close(unblock)
			}
			log.Printf("Content received from server %v", res.GetResult())
		}
		close(unblock)
	}()

	<-unblock
}
