package main

import (
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
}
