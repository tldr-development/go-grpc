package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/hojin-kr/fiber-grpc/gcp_datastore"
	"github.com/hojin-kr/fiber-grpc/proto"
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

func (s *server) Add(_ context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &proto.Response{Result: result}, nil
}

func (s *server) Multiply(_ context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a * b

	return &proto.Response{Result: result}, nil
}

func (s *server) DataStore(_ context.Context, request *proto.Request) (*proto.Response, error) {
	if env != "live" {
		log.Printf("daastore rpc")
	}
	key, value := request.GetA(), request.GetB()
	gcpDataStoreClient := gcp_datastore.GetClient(context.Background())
	kind := "Test"
	incompleteKey := gcp_datastore.IncompleteKey(kind, nil)
	completeKey, err := gcpDataStoreClient.Put(context.Background(), incompleteKey, &gcp_datastore.Test{Key: strconv.Itoa(int(key)), Value: strconv.Itoa(int(value))})
	if err != nil {
		return nil, err
	}
	result := completeKey.ID

	return &proto.Response{Result: result}, nil
}
