package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Println("Greet function called")
	fname := in.GetGreeting().GetFirstName()
	lname := in.GetGreeting().GetLastName()
	result := "Hello " + fname + " " + lname

	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func main() {
	fmt.Println("Hello world!")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		fmt.Println("Error!")
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error occured!")
	}

}
