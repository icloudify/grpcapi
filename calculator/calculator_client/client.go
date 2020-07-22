package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)
	doServerStreaming(c)

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

	for  {
		resStr, err := res.Recv()
		if err != nil {
			log.Fatal("Could not find result!", err)
		}

		fmt.Printf("Prime decomposition factor of %+v is %+v\n", req.Number, resStr.PrimeFactor)
	}
}
