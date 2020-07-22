package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct {
}

func (s *server) SquareRoot(ctx context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received request for square root....")

	if request.GetNumber() < 1 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %v\n", request.GetNumber()),
		)
	}
	sroot := math.Sqrt(float64(request.GetNumber()))
	resp := &calculatorpb.SquareRootResponse{
		NumberRoot: sroot,
	}

	return resp, nil
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Received FindMaximum....")
	max := int32(0)
	for {
		rec, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("********************\nFinal Max rceived %v\n********************\n", max)
			return stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
		}
		if max < rec.GetNumber() {
			max = rec.GetNumber()
			fmt.Println("New Max rceived ", max)
		}
	}

}

func (s *server) ComputeAverage(averageServer calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received ComputeAverage....")
	number := int32(0)
	count := int32(0)
	for {
		rec, err := averageServer.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			fmt.Printf("total %v and count %v\n", number, count)
			return averageServer.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float64(number / count),
			})
		}
		number = number + rec.GetNumber()
		count++
	}

	return nil
}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Calculator function called")
	fn := in.GetFirstNumber()
	sn := in.GetSecondNumber()
	result := fn + sn

	res := calculatorpb.SumResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Println("Received PrimeNumberDecomposition %v", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Println("Divisor has increaseed to%v", divisor)
		}
	}

	return nil
}

func main() {
	fmt.Println("***********************\n       GRPC Server       \n***********************\n")
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
