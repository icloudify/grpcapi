package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Calculator function called")
	fn := in.GetFirstNumber()
	sn := in.GetSecondNumber()
	result := fn+sn

	res := calculatorpb.SumResponse{
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
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error occured!")
	}
}
