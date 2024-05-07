package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/hojin-kr/go-grpc/gcp/datastore"
	proto "github.com/hojin-kr/go-grpc/wallet/proto"
	wallet_struct "github.com/hojin-kr/go-grpc/wallet/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")
var app = os.Getenv("APP")

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
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

// Get
func (s *server) Get(_ context.Context, request *proto.Request) (*proto.Response, error) {
	uuid := request.GetUuid()
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wallet")

	_wallet := &wallet_struct.Wallet{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, uuid, nil), _wallet)

	if _wallet.UUID == "" {
		_wallet.UUID = uuid
		_wallet.Status = "active"
		_wallet.Ticket = int64(10) // default ticket
		_wallet.Created = int64(time.Now().Unix())
		_wallet.Updated = int64(time.Now().Unix())
		dbClient.Put(context.Background(), datastore.NameKey(kind, uuid, nil), _wallet)
	}

	log.Printf("Get uuid: %v", uuid)
	return &proto.Response{Uuid: uuid, Status: _wallet.Status, Created: _wallet.Created, Updated: _wallet.Updated}, nil
}

// Update
func (s *server) Update(_ context.Context, request *proto.Request) (*proto.Response, error) {
	uuid := request.GetUuid()
	ticket := request.GetTicket()
	status := request.GetStatus()
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wallet")

	_wallet := &wallet_struct.Wallet{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, uuid, nil), _wallet)

	_wallet.Ticket = ticket
	_wallet.Status = status
	_wallet.Updated = int64(time.Now().Unix())

	dbClient.Put(context.Background(), datastore.NameKey(kind, uuid, nil), _wallet)

	log.Printf("Update uuid: %v, ticket: %v", uuid, ticket)
	return &proto.Response{Uuid: uuid, Status: _wallet.Status, Created: _wallet.Created, Updated: _wallet.Updated}, nil
}
