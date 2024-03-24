package main

import (
	"context"
	"net"
	"testing"

	proto "github.com/hojin-kr/fiber-grpc/inspire/proto"
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

func TestInspire(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{
		Prompt:  "안녕하세요, 저는 지금 약간은 우울하고 하는일이 잘 안되고 블루한 기분을 가지고 있어요.",
		Context: "심리 상담가의 입장에서 140자 정도의 편지를 적어주는데, 심리 상담가라는 것을 알리지 않고, 그냥 친구라고 생각하고 적어주세요. \n 그리고 이전 편지는 '안녕, 친구야. 최근에 기분이 안 좋다는 소식을 들었는데, 정말 미안해. 너를 걱정하게 돼. 힘든 시간을 보내고 있다는 걸 알아, 하지만 너가 홀로 아니라는 걸 잊지 마. 언제든 나에게 연락할 수 있어. 내가 도와 줄 수 있는 일이 없을지도 몰라도, 언제든지 귀 기울이고, 너를 위한 시간을 갖겠다는 건 알아줬으면 좋겠어. 앞으로 힘들면 나에게 알려줘. 항상 곁에 있을게.' 이렇게 적었으니까 참고해줘",
		Uuid:    "870ac164-9143-4165-aea1-1d93c3673e67",
	}
	res, err := client.Inspire(ctx, req)
	if err != nil {
		t.Fatalf("Inspire failed: %v", err)
	}
	if res == nil {
		t.Fatalf("Expected response, got %v", res)
	}
}

func TestSendNotifications(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)
	req := &proto.Request{}
	res, err := client.SendNotifications(ctx, req)
	if err != nil {
		t.Fatalf("SendNotifications failed: %v", err)
	}
	if res == nil {
		t.Fatalf("Expected response, got %v", res)
	}
}
