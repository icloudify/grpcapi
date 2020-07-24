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
	"time"
)

func main() {
	opts := grpc.WithInsecure()
	isTls := false
	//conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if isTls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatal("SSL authentication failed while loading certificate!", err)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	//doUnaryWithDeadline(c, 5*time.Second)
	//doUnaryWithDeadline(c, 1*time.Second)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	waitc := make(chan struct{})
	// Send a bunch of message
	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ravindra",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nath",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Thakur",
			},
		},
	}
	go func() {
		for _, num := range request {
			fmt.Println("Sending message ", num)
			stream.Send(num)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive a bunch of message
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Fatal("Could not receive from stream!", err)
			}
			if err == io.EOF {
				close(waitc)
			}

			fmt.Printf("Received %v\n", resp.GetResult())
		}
	}()
	// Block untill everything is done
	<-waitc
}
func doUnary(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ravindra",
			LastName:  "Thakur",
		},
	}
	res, err := c.Greet(context.Background(), &req)

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}
	fmt.Println(res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ravindra",
			LastName:  "Thakur",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), &req)

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			log.Println("EOF")
			break
		}

		fmt.Println(msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ravindra",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Thakur",
			},
		},
	}

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	for _, num := range request {
		fmt.Println("Sending number ", num)
		stream.Send(num)
		time.Sleep(1000 * time.Millisecond)
	}

	resStream, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	fmt.Println(resStream.GetResult())
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeouts time.Duration) {
	req := greetpb.GreetRequestWithDeadline{
		Greeting: &greetpb.Greeting{
			FirstName: "Ravindra",
			LastName:  "Thakur",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeouts)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, &req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline exceeded!")
			} else {
				fmt.Println("Unexpected error while calling deadline exceeded!")
			}
		} else {
			log.Fatal("Error while calling GreetWithDeadline ", err)
		}

		log.Fatal("Error while calling GreetWithDeadline ", err)
	}

	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
