package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/deepudoit/coolgo/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial(":50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect..%v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)
	// doServerStream(c)
	// doClientStream(c)
	// doBiDiStream(c)
	doDeadLine(c, 5*time.Second)
	doDeadLine(c, 1*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Pradeep",
			LastName:  "Gandla",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Pradeep",
			LastName:  "Gandla",
		},
	}
	resStream, err := c.GreeManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling stream: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something is not right: %v", err)
		}
		log.Printf("Response from server: %v", msg.GetResult())
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetReq{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Pradeep",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Hawking",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Diesel",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Sandhya",
			},
		},
	}

	msgStream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Failed to pass stream: %v", err)
	}

	for _, req := range requests {
		msgStream.Send(req)
		log.Println("Sending....")
		time.Sleep(time.Second * 1)
	}
	res, err := msgStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error in receiving stream: %v", err)
	}
	log.Printf("Server final msg: %v", res.Result)
}

func doBiDiStream(c greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetEveReq{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Apollo",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Endurance",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gargantua",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Do not go gentle into that good nite......",
			},
		},
	}
	msgStream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the connection: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Println("Sending....")
			msgStream.SendMsg(req)
			// time.Sleep(time.Millisecond * 200)
		}
		msgStream.CloseSend()
		log.Println("Sent all we got...Now waiting to get...")
	}()

	go func() {
		for {
			time.Sleep(time.Second * 3)
			msg, err := msgStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to get response from stream: %v", err)
				break
			}
			log.Printf("Hurray, receiving...%v\n", msg.Result)
		}
		close(waitc)
	}()

	<-waitc
}

func doDeadLine(c greetpb.GreetServiceClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &greetpb.GreetDeadlineReq{
		Greeting: &greetpb.Greeting{
			FirstName: "Pradeep",
			LastName:  "CloudDev",
		},
	}
	res, err := c.GreetDeadline(ctx, req)
	if err != nil {
		stsErr, ok := status.FromError(err)
		if ok {
			if stsErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout.....Deadline exceeded")
			} else {
				fmt.Printf("Another unexpected gRPC error: %v", stsErr)
			}
		} else {
			log.Fatalf("Failed to get response from server: %v", err)
		}
		return
	}
	log.Printf("Response :%v", res.Result)
}
