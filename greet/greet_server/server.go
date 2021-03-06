package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetRequestWithDeadline) (*greetpb.GreetResponseWithDeadline, error) {
	log.Println("GreetWithDeadline function called")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client Cancelled the request")
			return nil, status.Error(codes.DeadlineExceeded, "Client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}

	fname := req.GetGreeting().GetFirstName()
	lname := req.GetGreeting().GetLastName()
	result := "Hello " + fname + " " + lname

	res := greetpb.GreetResponseWithDeadline{
		Result: result,
	}
	return &res, nil
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
	fmt.Println("***********************\n       GRPC Server       \n***********************\n")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		fmt.Println("Error!")
	}

	opts := []grpc.ServerOption{}

	isTls := false

	if isTls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

		if err != nil {
			log.Fatal("SSL authentication failed while loading certificate!", err)
			return
		}

		opts = append(opts, grpc.Creds(creds))
	}
	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error occured!")
	}

}
