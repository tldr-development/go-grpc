package main

import (
	"context"
	"net"
	"testing"

	proto "github.com/hojin-kr/fiber-grpc/proto/sample"
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

func TesetAdd(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{A: 5, B: 3}
	res, err := client.Add(ctx, req)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if res.Result != 8 {
		t.Fatalf("Expected 8, got %d", res.Result)
	}
}

func TestMultiply(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{A: 5, B: 3}
	res, err := client.Multiply(ctx, req)
	if err != nil {
		t.Fatalf("Multiply failed: %v", err)
	}
	if res.Result != 15 {
		t.Fatalf("Expected 15, got %d", res.Result)
	}
}

func TestDatasotre(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{A: 5, B: 3}
	res, err := client.DataStore(ctx, req)
	if err != nil {
		t.Fatalf("DataStore failed: %v", err)
	}
	if res.Result < 1 {
		t.Fatalf("IncompleteKey None, got %d", res.Result)
	}
}

func TestCloudStorage(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.RequestSignedURL{Filename: "test", ContentType: "text/plain", Size: 1024}

	res, err := client.CloudStorage(ctx, req)
	if err != nil {
		t.Fatalf("CloudStorage failed: %v", err)
	}
	if res.Url == "" {
		t.Fatalf("SignedUrl is empty")
	}
}
