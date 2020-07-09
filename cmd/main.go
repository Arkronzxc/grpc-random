package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "test/internal/proto"
	"test/internal/service"
)

type Server struct {
	mt *service.MT
}

func (s *Server) GetRandomNumbers(_ context.Context, req *pb.RandomNumbersRequest) (*pb.RandomNumbersResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}
	numbers := s.mt.NextNAsync(req.Number, req.Max)

	return &pb.RandomNumbersResponse{Numbers: numbers}, nil
}

func (s *Server) Healtcheck(context.Context, *empty.Empty) (*pb.Ok, error) {
	return &pb.Ok{Response: "ok"}, nil
}

func (s *Server) validateRequest(req *pb.RandomNumbersRequest) error {
	if req.Number <= 0 || req.Max <= 0 {
		return errors.New("invalid request. Numbers and max must be positive")
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", "8080"))
	if err != nil {
		log.Fatalf("Error listening %v", err)
	}

	grpcServer := grpc.NewServer()
	s := &Server{}

	s.mt = service.New(1234)
	pb.RegisterRandomNumberGeneratorServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Can not serve port 8080 with error %v", err)
	}
}
