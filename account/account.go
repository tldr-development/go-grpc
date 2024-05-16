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
	token := request.GetToken()
	platform := request.GetPlatform()

	if token == "" || platform == "" {
		log.Printf("Failed New Account: %v", accountUUID)
		return &proto.Response{Uuid: ""}, nil
	}
	// 1. platform에서 조회
	kind := datastore.GetKindByPrefix(app+":"+env, "platform")
	newPlatform := account.Platform{
		AccountID: accountUUID,
		Token:     token,
		Platform:  platform,
	}

	platforms := []account.Platform{}

	dbClient := datastore.GetClient(context.Background())
	query := datastore.NewQuery(kind).FilterField("token", "=", token).FilterField("platform", "=", platform).Limit(1)
	dbClient.GetAll(context.Background(), query, &platforms)

	if len(platforms) > 0 {
		log.Printf("platforms: %v", platforms)
		kind = datastore.GetKindByPrefix(app+":"+env, "account")
		account := &account.Account{}
		dbClient.Get(context.Background(), datastore.NameKey(kind, platforms[0].AccountID, nil), account)
		return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
	}

	// 2. platform이 없으면 새로 생성
	// 2-1. account 생성
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

	// 2-2. platform 생성
	_, err = dbClient.Put(context.Background(), datastore.IncompleteKey(kind, nil), &newPlatform)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}
	// 3. account에서 조회
	log.Printf("account_uuid: %s", accountUUID)
	account := &account.Account{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), account)
	return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
}

func (s *server) Add(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	token := request.GetToken()
	platform := request.GetPlatform()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "platform")

	newPlatform := account.Platform{
		AccountID: accountUUID,
		Token:     token,
		Platform:  platform,
	}

	_, err := dbClient.Put(context.Background(), datastore.IncompleteKey(kind, nil), &newPlatform)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("newPlatform: %v", newPlatform)

	account := &account.Account{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), account)
	return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
}
