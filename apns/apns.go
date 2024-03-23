package apns

import (
	"context"
	"log"
	"net"
	"os"

	proto "github.com/hojin-kr/fiber-grpc/account/proto"
	"github.com/hojin-kr/fiber-grpc/apns"
	"github.com/hojin-kr/fiber-grpc/gcp/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

type Apns struct {
	UUID  string // Account uuid
	Token string // apns
}

var env = os.Getenv("ENV")
var app = os.Getenv("APP")

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

// set apns token using uuid
func (s *server) SetToken(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	apnsToken := request.GetToken()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+env, "apns")

	apns := &apns.Apns{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), apns)

	apns.UUID = accountUUID
	apns.Token = apnsToken

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, accountUUID, nil), apns)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("apns: %v", apns)
	return &proto.Response{Uuid: apns.uuid, Token: apns.token}, nil
}
