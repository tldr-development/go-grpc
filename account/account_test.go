package main

import (
	"context"
	"net"
	"testing"

	proto "github.com/hojin-kr/go-grpc/account/proto"
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

func TestInit(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	// 기존 계정 조회
	req := &proto.Request{
		Uuid:     "",
		Token:    "000166.bde014069c2b4c0994bfeeaa490cfc39.1249",
		Platform: "apple",
	}
	res, err := client.Init(ctx, req)
	t.Logf("res: %v", res)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if res.Uuid == "" {
		t.Fatalf("Expected uuid, got %s", res.Uuid)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{
		Uuid:     "53954449-0089-4558-9509-c5734e4d79ba",
		Token:    "c5c2c0f5e073e434497af469e1f3c66b9.0.srww.kDLHVRwQxg1OUEhN18-5gA",
		Platform: "apple",
	}
	res, err := client.Delete(ctx, req)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if res.Uuid != "" {
		t.Fatalf("Expected empty uuid, got %s", res.Uuid)
	}
}
