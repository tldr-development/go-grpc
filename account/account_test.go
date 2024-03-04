package main

import (
	"context"
	"net"
	"testing"

	proto "github.com/hojin-kr/fiber-grpc/account/proto"
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
	req := &proto.Request{Uuid: "3364c611-0f37-4f3d-bf96-295dc8d3c56e"}
	res, err := client.Init(ctx, req)
	t.Logf("res: %v", res)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if res.Uuid == "" {
		t.Fatalf("Expected uuid, got %s", res.Uuid)
	}

	// 신규 계정 생성
	req = &proto.Request{}
	res, err = client.Init(ctx, req)
	t.Logf("res: %v", res)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if res.Uuid == "" {
		t.Fatalf("Expected uuid, got %s", res.Uuid)
	}

}
