package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)

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
