package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/edganiukov/apns"
	proto "github.com/hojin-kr/go-grpc/apns/proto"
	"github.com/hojin-kr/go-grpc/gcp/datastore"
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

var (
	apple_team_id     = os.Getenv("APPLE_TEAM_ID")
	apple_bundle_id   = os.Getenv("APPLE_BUNDLE_ID")
	apple_apns_key_id = os.Getenv("APPLE_APNS_KEY_ID")
	apple_apns_key    = os.Getenv("APPLE_APNS_KEY")
)

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

// set apns token using uuid
func (s *server) SetToken(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	apnsToken := request.GetToken()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "apns")

	_apns := &Apns{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), _apns)

	_apns.UUID = accountUUID
	_apns.Token = apnsToken

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, accountUUID, nil), _apns)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("apns: %v", _apns)
	return &proto.Response{Uuid: _apns.UUID, Token: _apns.Token}, nil
}

// get apns token using uuid
func (s *server) GetToken(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "apns")

	_apns := &Apns{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), _apns)

	log.Printf("apns: %v", _apns)
	return &proto.Response{Uuid: _apns.UUID, Token: _apns.Token}, nil
}

// send notification using apns token
func (s *server) SendNotification(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	title := request.GetTitle()
	subtitle := request.GetSubtitle()
	body := request.GetBody()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "apns")

	_apns := &Apns{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), _apns)
	if _apns.Token == "" {
		log.Printf("No token")
		return &proto.Response{}, nil
	}

	notification([]string{_apns.Token}, title, subtitle, body)

	return &proto.Response{Uuid: _apns.UUID, Token: _apns.Token}, nil
}

func notification(apnsTokens []string, title string, subtitle string, body string) {
	const (
		DevelopmentGateway = "https://api.sandbox.push.apple.com"
		ProductionGateway  = "https://api.push.apple.com"
	)
	GateWay := DevelopmentGateway
	if env == "prod" {
		GateWay = ProductionGateway
	}
	_apple_apns_key, _ := b64.StdEncoding.DecodeString(apple_apns_key)
	c, err := apns.NewClient(
		apns.WithJWT(_apple_apns_key, apple_apns_key_id, apple_team_id),
		apns.WithBundleID(apple_bundle_id),
		apns.WithMaxIdleConnections(10),
		apns.WithTimeout(5*time.Second),
		apns.WithEndpoint(GateWay),
	)
	if err != nil {
		print(err)
		/* ... */
	}
	for i := 0; i < len(apnsTokens); i++ {
		_, err := c.Send(apnsTokens[i],
			apns.Payload{
				APS: apns.APS{
					Alert: apns.Alert{
						Title:    title,
						Subtitle: subtitle,
						Body:     body,
					},
					Sound: "default",
				},
			},
			apns.WithExpiration(10),
			apns.WithPriority(5),
		)
		if err != nil {
			print(err)
			/* ... */
		}
	}
}
