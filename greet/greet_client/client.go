package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/deepudoit/coolgo/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
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
	doClientStream(c)
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
