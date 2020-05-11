package main

import (
	"context"
	"log"
	"net"

	"github.com/deepudoit/coolgo/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Sum(ctx context.Context, r *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	result := r.FirstNum + r.SecNum
	res := &calcpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to connect 500551: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
