package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/deepudoit/coolgo/gogrpc/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
	//doClientStream(c)
	doSqrt(c)
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

func doSqrt(c calcpb.CalculatorServiceClient) {
	//Correct call
	doErrSqrt(c, -10)

	//Error call
	doErrSqrt(c, -4)
}

func doErrSqrt(c calcpb.CalculatorServiceClient, num int32) {
	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootReq{Num: num})

	if err != nil {
		if gErr, ok := status.FromError(err); ok {
			fmt.Printf("%v\n", gErr.Message())
			return
		}
		log.Fatalf("Falied to call remote method: %v", err)
		return
	}
	log.Printf("Here's the sqrt of %d: %.2f\n", num, res.GetResult())
}
