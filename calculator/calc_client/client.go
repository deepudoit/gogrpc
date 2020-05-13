package main

import (
	"context"
	"io"
	"log"
	"time"

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

	// doUnary(c)
	// doServerStream(c)
	doClientStream(c)
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

func doServerStream(c calcpb.CalculatorServiceClient) {
	req := &calcpb.PrimeNumDecomReq{
		Num: 500,
	}
	res, err := c.PrimeNumDecom(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed unable to call server method: %v", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to get the server response: %v", err)
		}
		log.Printf("Prime decomposition factors: %v", msg)
	}
}

func doClientStream(c calcpb.CalculatorServiceClient) {
	nums := []*calcpb.ComputeAvgReq{
		{
			Num: 10,
		},
		{
			Num: 30,
		},
		{
			Num: 50,
		},
		{
			Num: 44,
		},
		{
			Num: 55,
		},
		{
			Num: 10,
		},
		{
			Num: 77,
		},
	}

	resStream, err := c.ComputeAvg(context.Background())
	if err != nil {
		log.Fatalf("Failed to get server response: %v", err)
	}

	for _, num := range nums {
		resStream.Send(num)
		log.Println("Sending....")
		time.Sleep(time.Second * 1)
	}

	res, err := resStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error in receiving from stream: %v", err)
	}
	log.Printf("Average : %.2f", res.Avg)
}
