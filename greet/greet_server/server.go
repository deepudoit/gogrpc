package main

import (
	"log"
	"net"

	"github.com/deepudoit/coolgo/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to connect..%v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server {})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve..%v", err)
	}
}

