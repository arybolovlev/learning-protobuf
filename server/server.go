package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/arybolovlev/learning-protobuf/proto"
)

var (
	port = flag.Int("port", 1602, "The server port")
)

type server struct {
	pb.NumberServer
}

func (s *server) IncreseNumber(ctx context.Context, n *pb.NumberRequest) (*pb.NumberResponse, error) {
	log.Printf("Received number: %v", n.Number)
	return &pb.NumberResponse{
		Number:    n.Number + 1,
		OldNumber: n.Number,
	}, nil
}

func Run() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Falied to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNumberServer(s, &server{})
	log.Printf("server listening at %v", l.Addr())

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Println("Bye!")
}
