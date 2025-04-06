package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	proto "github.com/hojin-kr/go-grpc/wscan/proto"
	"google.golang.org/grpc"
)

func main() {
	// gRPC 서버 주소
	addr := "localhost:50051"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAddServiceClient(conn)

	// 이미지 파일 읽기
	imagePath := "sample.jpg" // 테스트할 이미지 경로
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}

	// 요청 생성
	req := &proto.Request{
		Image:   imageBytes,
		Context: "My suitcase for a 3-day business trip.",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Wscan 호출
	resp, err := client.Wscan(ctx, req)
	if err != nil {
		log.Fatalf("could not call Wscan: %v", err)
	}

	fmt.Println("UUID:", resp.GetUuid())
	fmt.Println("Prompt:", resp.GetPrompt())
	fmt.Println("Message:", resp.GetMessage())
}
