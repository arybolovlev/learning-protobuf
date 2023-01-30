package main

import (
	"context"
	"flag"
	"log"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	i, err := c.IncreseNumber(ctx, &pb.NumberRequest{Number: 4})
	if err != nil {
		log.Fatalf("could not increase the number: %v", err)
	}
	log.Printf("New number value: %v", i.Number)
	log.Printf("Old number value: %v", i.OldNumber)
}
