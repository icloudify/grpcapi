package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with \n")

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while receiveing stream: %v", err)
			return err
		}

		fname := res.GetGreeting().GetFirstName()
		result := "hello " + fname + "!"
		error := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if error != nil {
			log.Fatalf("Error while sending on stream: %v", err)
			return error
		}

	}
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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req.GetGreeting().GetFirstName())
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManytimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "!"
	}
}

func main() {
	fmt.Println("Hello world!")
	lis, err := net.Listen("tcp", "0.0.0.0:50052")

	if err != nil {
		fmt.Println("Error!")
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error occured!")
	}

}
