package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/google/uuid"
	apns_proto "github.com/hojin-kr/go-grpc/apns/proto"
	"github.com/hojin-kr/go-grpc/gcp/datastore"
	proto "github.com/hojin-kr/go-grpc/inspire/proto"
	inspire_struct "github.com/hojin-kr/go-grpc/inspire/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")
var app = os.Getenv("APP")
var projectID = os.Getenv("PROJECT_ID")
var apns_server = os.Getenv("APNS_SERVER")

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
	prompt := request.GetPrompt()
	gen_context := request.GetContext()
	_uuid := request.GetUuid()

	generateByGemini(prompt, gen_context, _uuid)

	return &proto.Response{}, nil
}

// 내 inspire 목록을 조회
func (s *server) GetInspires(_ context.Context, request *proto.Request) (*proto.Responses, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "complete").Order("-Created").Limit(100)
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	responses := []*proto.Response{}
	for _, inspire := range inspires {
		responses = append(responses, &proto.Response{Uuid: inspire.UUID, Prompt: inspire.Prompt, Message: inspire.Message, Created: inspire.Created, Updated: inspire.Updated})
		log.Printf("inspire: %v", inspire)
	}

	// return inspires list to client
	return &proto.Responses{Responses: responses}, nil
}

// 내 마지막 inspire를 조회
func (s *server) GetLastInspire(_ context.Context, request *proto.Request) (*proto.Response, error) {
	inspires := getLastInspire(request.GetUuid())
	inspire := inspire_struct.Inspire{}
	for _, inspire = range inspires {
		break
	}
	response := &proto.Response{Uuid: inspire.UUID, Prompt: inspire.Prompt, Message: inspire.Message, Created: inspire.Created, Updated: inspire.Updated}
	log.Printf("inspire: %v", inspire)

	// return inspire to client
	return response, nil
}

// 특정 시간 이후의 inspire last를 조회해서 inspire를 생성한다.
func (s *server) GenerateInspireAfterCreatedLast(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "complete").FilterField("Updated", "<", request.GetCreated()).DistinctOn("UUID").Limit(100)
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	for _, inspire := range inspires {
		go generateByGemini(inspire.Prompt, inspire.Context, inspire.UUID)
		log.Println("inspire: ", inspire)
	}

	return &proto.Response{}, nil
}

// get Last inspire by last one
func getLastInspire(_uuid string) []inspire_struct.Inspire {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", _uuid).FilterField("Status", "=", "complete").Order("-Created").Limit(1)
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	return inspires
}

// 내 inspire를 삭제
func (s *server) DeleteInspire(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	inspire := &inspire_struct.Inspire{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), inspire)

	inspire.Status = "deleted"
	inspire.Updated = int64(time.Now().Unix())

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), inspire)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}

	log.Printf("inspire: %v", inspire)

	// return inspire to client
	return &proto.Response{Uuid: inspire.UUID, Prompt: inspire.Prompt, Message: inspire.Message, Created: inspire.Created, Updated: inspire.Updated}, nil
}

// SendNotification 특정 유저의 inspire를 조회하여 pending 상태만 notification을 보낸다.
func (s *server) SendNotification(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "pending")
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	wg := sync.WaitGroup{}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(apns_server, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	c := apns_proto.NewAddServiceClient(conn)

	for _, inspire := range inspires {
		if inspire.NameKey == "" {
			log.Print("continue")
			continue
		}
		// request to notification grpc server
		wg.Add(1)
		go invokeNotification(c, inspire, &wg)
		// inspire의 status를 complete로 변경
		inspire.Status = "complete"
		inspire.Updated = int64(time.Now().Unix())

		_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, inspire.NameKey, nil), &inspire)
		if err != nil {
			log.Printf("Failed to put: %v", err)
		}
		log.Print("inspire notification : ", inspire.NameKey)
	}
	wg.Wait()

	return &proto.Response{}, nil
}

func (s *server) SendNotifications(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// pendding 상태의 inspire를 조회하여 notification을 보낸다.
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "pending").DistinctOn("UUID")
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	wg := sync.WaitGroup{}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(apns_server, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	c := apns_proto.NewAddServiceClient(conn)

	for _, inspire := range inspires {
		if inspire.NameKey == "" {
			log.Print("continue")
			continue
		}
		// request to notification grpc server
		wg.Add(1)
		go invokeNotification(c, inspire, &wg)
		// inspire의 status를 complete로 변경
		inspire.Status = "complete"
		inspire.Updated = int64(time.Now().Unix())

		_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, inspire.NameKey, nil), &inspire)
		if err != nil {
			log.Printf("Failed to put: %v", err)
		}
		log.Print("inspires notification : ", inspire.NameKey)
	}
	wg.Wait()

	return &proto.Response{}, nil
}

func generateByGemini(prompt, gen_context, _uuid string) []string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(model)
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt+"\n"+gen_context))
	if err != nil {
		log.Fatal(err)
	}

	parts := printResponse(resp)

	// set last inspire
	for _, part := range parts {
		fmt.Println(part + "\n")
		setInpireDatastore(_uuid, prompt, gen_context, part)
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

func setInpireDatastore(_uuid, prompt, gen_context, message string) string {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	inspire := &inspire_struct.Inspire{}
	inspire.UUID = _uuid
	inspire.Prompt = prompt
	inspire.Context = gen_context
	inspire.Message = message
	inspire.Created = int64(time.Now().Unix())
	inspire.Status = "pending"
	inspire.NameKey = uuid.New().String()

	NameKey := datastore.NameKey(kind, inspire.NameKey, nil)
	_, err := dbClient.Put(context.Background(), NameKey, inspire)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}
	log.Printf("inspire: %v", inspire)
	return inspire.NameKey
}

func invokeNotification(c apns_proto.AddServiceClient, inspire inspire_struct.Inspire, wg *sync.WaitGroup) {
	ctx := context.Background()
	_, err := c.SendNotification(ctx, &apns_proto.Request{Uuid: inspire.UUID, Title: "", Subtitle: "", Body: inspire.Message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	wg.Done()
}
