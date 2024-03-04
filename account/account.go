package main

import (
	"context"
	"log"
	"net"
	"os"

	proto "github.com/hojin-kr/fiber-grpc/account/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")

func main() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	if env != "live" {
		log.Printf("Run server")
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(err)
	}
}

func (s *server) Init(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// Sign Up if Request uuid is empty

	return &proto.Response{Result: 0}, nil
}
