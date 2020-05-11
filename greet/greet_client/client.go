package main

import (
	"context"
	"log"

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
	doUnary(c)
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
