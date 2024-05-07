package main

import (
	"context"
	"log"
	"net"
	"testing"

	proto "github.com/hojin-kr/go-grpc/wallet/proto"
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

func TestGet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{Uuid: "3364c611-0f37-4f3d-bf96-295dc8d3c56a"}
	res, err := client.Get(ctx, req)
	t.Logf("res: %v", res)
	log.Println(res)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if res.Uuid == "" {
		t.Fatalf("Expected uuid, got %s", res.Uuid)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{Uuid: "3364c611-0f37-4f3d-bf96-295dc8d3c56a", Ticket: 5}
	res, err := client.Update(ctx, req)
	t.Logf("res: %v", res)
	log.Println(res)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if res.Uuid == "" {
		t.Fatalf("Expected uuid, got %s", res.Uuid)
	}
}
