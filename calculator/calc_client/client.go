package main

import (
	"context"
	"log"

	"github.com/deepudoit/coolgo/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial(":50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Faile to connect 50051: %v", err)
	}

	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)

	doUnary(c)
}

func doUnary(c calcpb.CalculatorServiceClient) {
	req := &calcpb.SumRequest{
		FirstNum: 10,
		SecNum:   20,
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("No response from server: %v", err)
	}

	log.Printf("Sum : %d\n", res.SumResult)
}
