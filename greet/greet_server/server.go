package main

import (
	"log"
	"net"
	"context"

	"github.com/deepudoit/coolgo/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	
}

func (s *server) Greet(ctx context.Context, r *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := r.GetGreeting().GetFirstName()
	lastName := r.GetGreeting().GetLastName()
	res := "Welcome... " + firstName + ", " + lastName

	response := &greetpb.GreetResponse{
		Result: res,
	} 

	return response, nil
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