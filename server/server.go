package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "github.com/arybolovlev/learning-protobuf/proto"
)

var (
	port = flag.Int("port", 1602, "The server port")
)

type server struct {
	pb.NumberServer
	subscribers sync.Map
}

type subscriber struct {
	eventId  int32
	stream   pb.Number_SubscribeServer
	finished chan<- bool
}

func (s *server) Subscribe(req *pb.SubscribeRequest, stream pb.Number_SubscribeServer) error {
	log.Println("New subscribe request received")

	f := make(chan bool)
	id := uuid.New()

	s.subscribers.Store(id.String(), subscriber{eventId: req.EventId, stream: stream, finished: f})

	log.Printf("New subscribe request ID %v with event ID: %v", id, req.EventId)

	ctx := stream.Context()

	for {
		select {
		case <-f:
			log.Println("Closing streaming for client")
			return nil
		case <-ctx.Done():
			log.Println("Client has disconnected")
			return nil
		}
	}
}

func (s *server) Unsubscribe(ctx context.Context, us *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	_, ok := s.subscribers.LoadAndDelete(us.Id)
	if !ok {
		log.Printf("Failed to unsubsribe ID %v", us.Id)
		return &pb.UnsubscribeResponse{Unsubscribed: false}, fmt.Errorf("Failed to unsubsribe ID %v", us.Id)
	}
	log.Printf("Unsubsribed ID %v", us.Id)
	return &pb.UnsubscribeResponse{Unsubscribed: true}, nil
}

func (s *server) GenerateData() {
	log.Println("Starting data generation")

	for {
		time.Sleep(time.Second)

		s.subscribers.Range(func(k, v interface{}) bool {
			id, ok := k.(string)
			if !ok {
				log.Printf("Failed to cast subscriber key: %v", id)
				return false
			}
			sub, ok := v.(subscriber)
			if !ok {
				log.Printf("Failed to cast subscriber value: %T", v)
				return false
			}

			err := sub.stream.Send(&pb.SubscribeResponse{Id: id, EventId: sub.eventId})
			if err != nil {
				log.Println("Failed to send data to client")
			}
			sub.eventId++
			s.subscribers.Store(id, sub)

			return true
		})
	}
}

func Run() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Falied to listen: %v", err)
	}

	srv := &server{}

	go srv.GenerateData()

	s := grpc.NewServer()
	pb.RegisterNumberServer(s, srv)
	log.Printf("server listening at %v", l.Addr())

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Println("Bye!")
}
