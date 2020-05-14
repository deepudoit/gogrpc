package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/deepudoit/coolgo/gogrpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
}

func (*server) Greet(ctx context.Context, r *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := r.GetGreeting().GetFirstName()
	lastName := r.GetGreeting().GetLastName()
	res := "Welcome... " + firstName + ", " + lastName

	response := &greetpb.GreetResponse{
		Result: res,
	}

	return response, nil
}

func (*server) GreeManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreeManyTimesServer) error {
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 20; i++ {
		res := "Hello " + firstName + " You rock.." + strconv.Itoa(i)
		result := &greetpb.GreetManyTimesResponse{
			Result: res,
		}
		stream.Send(result)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Finished reading the msg")
			return stream.SendAndClose(&greetpb.LongGreetRes{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Failed to client stream messages: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += ">>>>" + firstName + "......\n"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to connect to stream: %v", err)
			return err
		}
		time.Sleep(time.Second * 2)
		log.Println("From Voyager...")
		firstname := req.GetGreeting().GetFirstName()
		result += ">>>>" + firstname + ".......\n"
		err = stream.Send(&greetpb.GreetEveRes{
			Result: firstname,
		})
		if err != nil {
			log.Fatalf("Failed to send response: %v", err)
		}
	}
}

func (*server) GreetDeadline(ctx context.Context, r *greetpb.GreetDeadlineReq) (*greetpb.GreetDeadlineRes, error) {
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if ctx.Err() == context.Canceled {
			fmt.Println("Client cancelled the request")
			return nil, status.Errorf(codes.Canceled, "Request was cancelled by client")
		}
	}
	firstName := r.GetGreeting().GetFirstName()
	lastName := r.GetGreeting().GetLastName()
	res := "Welcome... " + firstName + ", " + lastName

	response := &greetpb.GreetDeadlineRes{
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
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve..%v", err)
	}
}
