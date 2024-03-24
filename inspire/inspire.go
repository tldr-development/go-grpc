package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/hojin-kr/fiber-grpc/gcp/datastore"
	proto "github.com/hojin-kr/fiber-grpc/inspire/proto"
	inspire_struct "github.com/hojin-kr/fiber-grpc/inspire/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")
var app = os.Getenv("APP")
var projectID = os.Getenv("PROJECT_ID")

const location = "us-central1"
const model = "gemini-1.0-pro-001"

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

func (s *server) Inspire(_ context.Context, request *proto.Request) (*proto.Response, error) {
	prompt := request.GetPrompt() + "\n" + request.GetContext()
	uuid := request.GetUuid()

	generateByGemini(prompt, uuid)

	// requset to notification grpc server

	return &proto.Response{}, nil
}

func (s *server) SendNotifications(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// pendding 상태의 inspire를 조회하여 notification을 보낸다.
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("status", "=", "pending")
	inspires := []inspire_struct.Inspire{}
	keys, _ := dbClient.GetAll(context.Background(), query, &inspires)

	for _, key := range keys {
		// request to notification grpc server
		// inspire의 status를 complete로 변경
		inspires[key.ID].Status = "complete"
		_, err := dbClient.Put(context.Background(), datastore.IDKey(kind, key.ID, nil), inspires[key.ID])
		if err != nil {
			log.Printf("Failed to put: %v", err)
		}
	}

	return &proto.Response{}, nil
}

func generateByGemini(prompt string, uuid string) []string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(model)
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	parts := printResponse(resp)

	for _, part := range parts {
		fmt.Println(part + "\n")
		setInpireDatastore(uuid, prompt, part)
	}
	return parts
}

func printResponse(resp *genai.GenerateContentResponse) []string {
	var parts []string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			parts = append(parts, fmt.Sprint(part))
		}
	}
	return parts
}

func setInpireDatastore(uuid, prompt, message string) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+env, "inspire")

	inspire := &inspire_struct.Inspire{}
	inspire.UUID = uuid
	inspire.Prompt = prompt
	inspire.Message = message
	inspire.Created = strconv.Itoa(int(time.Now().Unix()))
	inspire.Status = "pending"

	_, err := dbClient.Put(context.Background(), datastore.IncompleteKey(kind, nil), inspire)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}
	log.Printf("inspire: %v", inspire)
}
