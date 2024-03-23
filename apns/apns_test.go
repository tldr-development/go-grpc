package main

import (
	"context"
	"net"
	"testing"

	proto "github.com/hojin-kr/fiber-grpc/apns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const buffSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(buffSize)
	s := grpc.NewServer()
	proto.RegisterAddServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSetToken(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{Uuid: "test-uuid", Token: "785568d447412ab9655333f1206e30275a18fba6a6b49a174652391cd3f9d009"}
	res, err := client.SetToken(ctx, req)
	if err != nil {
		t.Fatalf("SetToken failed: %v", err)
	}
	if res.Token == "" {
		t.Fatalf("Expected token, got %s", res.Token)
	}
}

func TestGetToken(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{Uuid: "test-uuid"}
	res, err := client.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if res.Token == "" {
		t.Fatalf("Expected token, got %s", res.Token)
	}
	t.Log(res)
}

func TestNotification(t *testing.T) {
	apnsTokens := []string{"785568d447412ab9655333f1206e30275a18fba6a6b49a174652391cd3f9d009", "785568d447412ab9655333f1206e30275a18fba6a6b49a174652391cd3f9d009"}
	title := "title"
	subtitle := "subtitle"
	body := "body"
	notification(apnsTokens, title, subtitle, body)

	t.Log("Notification success")
}

func TestSendNotification(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{Uuid: "test-uuid", Token: "785568d447412ab9655333f1206e30275a18fba6a6b49a174652391cd3f9d009"}
	res, err := client.SendNotification(ctx, req)
	if err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}
	if res.Token == "" {
		t.Fatalf("Expected token, got %s", res.Token)
	}
	t.Log(res)
}
