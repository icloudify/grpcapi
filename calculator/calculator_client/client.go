package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)

}
func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximum BiDi Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	request := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 11,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 5,
		},
	}

	waitc := make(chan struct{})

	// send go routine
	go func() {
		for _, number := range request {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(number)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of...: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.SumRequest{
		FirstNumber:  11,
		SecondNumber: 12,
	}
	res, err := c.Sum(context.Background(), &req)

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}
	fmt.Printf("Sum of %+v and %+v is %+v\n", req.FirstNumber, req.SecondNumber, res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12345678,
	}
	res, err := c.PrimeNumberDecomposition(context.Background(), &req)

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	for {
		resStr, err := res.Recv()
		if err != nil {
			log.Fatal("Could not find result!", err)
		}

		fmt.Printf("Prime decomposition factor of %+v is %+v\n", req.Number, resStr.PrimeFactor)
	}
}
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Average function was invoked with a streaming request\n")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	request := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 5,
		},
	}

	for _, num := range request {
		fmt.Println("Sending number ", num)
		stream.Send(num)
		time.Sleep(1000 * time.Millisecond)
	}

	avg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	fmt.Printf("Average %v", avg.GetAverage())

}
