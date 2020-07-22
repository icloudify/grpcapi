package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Received FindMaximum....")
	max := int32(0)
	for {
		rec, err := stream.Recv()

		if err == io.EOF {
			return stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
		}
		if max < rec.GetNumber() {
			max = rec.GetNumber()
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
