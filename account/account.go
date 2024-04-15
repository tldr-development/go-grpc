package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	proto "github.com/hojin-kr/go-grpc/account/proto"
	account "github.com/hojin-kr/go-grpc/account/struct"
	"github.com/hojin-kr/go-grpc/gcp/datastore"
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

func (s *server) Init(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// Sign Up if Request uuid is empty
	accountUUID := request.GetUuid()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+env, "account")

	if accountUUID != "" {
		log.Printf("account_uuid: %s", accountUUID)
		// DB에서 uuid로 조회
		account := &account.Account{}
		dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), account)
		return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
	}
	// Sign Up
	accountUUID = uuid.New().String()
	log.Printf("account_uuid: %s", accountUUID)

	timestampStr := strconv.Itoa(int(time.Now().Unix()))

	newAccount := account.Account{
		UUID:    accountUUID,
		Status:  "active",
		Created: timestampStr,
		Updated: timestampStr,
	}

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, accountUUID, nil), &newAccount)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("newAccount: %v", newAccount)
	return &proto.Response{Uuid: newAccount.UUID, Status: newAccount.Status, Created: newAccount.Created, Updated: newAccount.Updated}, nil
}
