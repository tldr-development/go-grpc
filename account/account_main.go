// account/account_main.go
package main

import (
	"log"
	"net"
	"os"

	"account/model"
	"account/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
	proto.RegisterAddServiceServer(srv, &model.Server{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(e)
	}
}
