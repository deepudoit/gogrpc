package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/deepudoit/coolgo/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (*server) PrimeNumDecom(req *calcpb.PrimeNumDecomReq, stream calcpb.CalculatorService_PrimeNumDecomServer) error {
	num := req.GetNum()
	divisor := int64(2)
	for num > 1 {
		if num%divisor == 0 {
			stream.Send(&calcpb.PrimNumDecomRes{
				PrimeFactor: divisor,
			})
			num /= divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased: %v\n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAvg(stream calcpb.CalculatorService_ComputeAvgServer) error {
	var sum int32 = 0
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Finished reading server stream")
			return stream.SendAndClose(&calcpb.ComputeAvgRes{
				Avg: float64(sum) / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("Failed to server stream requests: %v", err)
		}
		sum += req.GetNum()
		count++
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootReq) (*calcpb.SquareRootRes, error) {
	num := req.GetNum()
	log.Printf("Num from client: %v", num)
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("GRPC_ERR00012: Invalid num %v", num))
	}
	return &calcpb.SquareRootRes{
		Result: math.Sqrt(float64(num)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to connect 500551: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
