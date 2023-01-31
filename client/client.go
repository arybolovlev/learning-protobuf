package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/arybolovlev/learning-protobuf/proto"
)

var (
	addr = flag.String("addr", "localhost:1602", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNumberClient(conn)

	rand.Seed(time.Now().UnixNano())
	req := &pb.SubscribeRequest{EventId: int32(rand.Int())}
	serverStream, err := c.Subscribe(context.Background(), req)
	if err != nil {
		log.Fatal("could not stream")
	}

	counter := 0

	for {
		resp, err := serverStream.Recv()
		if err == io.EOF {
			break
		}

		log.Printf("Response [%v] - %v\n", resp.Id, resp.EventId)

		counter++
		if counter == 5 {
			r, err := c.Unsubscribe(context.Background(), &pb.UnsubscribeRequest{Id: resp.Id})
			if err != nil {
				log.Printf("Unsubscribe response error: %v\n", err)
			}
			if r.Unsubscribed {
				log.Printf("Successfully unsubscribed client ID %v", resp.Id)
				break
			}
		}
	}
}
