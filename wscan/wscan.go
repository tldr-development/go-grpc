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
	proto "github.com/hojin-kr/go-grpc/wscan/proto"
	wscan_struct "github.com/hojin-kr/go-grpc/wscan/struct"
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

func (s *server) Wscan(_ context.Context, request *proto.Request) (*proto.Response, error) {
	prompt := request.GetPrompt()
	gen_context := request.GetContext()
	_uuid := request.GetUuid()

	generateByGemini(prompt, gen_context, _uuid)

	return &proto.Response{}, nil
}

// 내 wscan 목록을 조회
func (s *server) GetWscans(_ context.Context, request *proto.Request) (*proto.Responses, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "complete").Order("-Created").Limit(100)
	wscans := []wscan_struct.Wscan{}
	dbClient.GetAll(context.Background(), query, &wscans)

	responses := []*proto.Response{}
	for _, wscan := range wscans {
		responses = append(responses, &proto.Response{Uuid: wscan.UUID, Prompt: wscan.Prompt, Message: wscan.Message, Created: wscan.Created, Updated: wscan.Updated})
		log.Printf("wscan: %v", wscan)
	}

	// return wscans list to client
	return &proto.Responses{Responses: responses}, nil
}

// 내 마지막 wscan를 조회
func (s *server) GetLastWscan(_ context.Context, request *proto.Request) (*proto.Response, error) {
	wscans := getLastWscan(request.GetUuid())
	wscan := wscan_struct.Wscan{}
	for _, wscan = range wscans {
		break
	}
	response := &proto.Response{Uuid: wscan.UUID, Prompt: wscan.Prompt, Message: wscan.Message, Created: wscan.Created, Updated: wscan.Updated}
	log.Printf("wscan: %v", wscan)

	// return wscan to client
	return response, nil
}

// 특정 시간 이후의 wscan last를 조회해서 wscan를 생성한다.
func (s *server) GenerateWscanAfterCreatedLast(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "complete").FilterField("Updated", "<", request.GetCreated()).DistinctOn("UUID").Limit(100)
	wscans := []wscan_struct.Wscan{}
	dbClient.GetAll(context.Background(), query, &wscans)

	for _, wscan := range wscans {
		go generateByGemini(wscan.Prompt, wscan.Context, wscan.UUID)
		log.Println("wscan: ", wscan)
	}

	return &proto.Response{}, nil
}

// get Last wscan by last one
func getLastWscan(_uuid string) []wscan_struct.Wscan {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", _uuid).FilterField("Status", "=", "complete").Order("-Created").Limit(1)
	wscans := []wscan_struct.Wscan{}
	dbClient.GetAll(context.Background(), query, &wscans)

	return wscans
}

// 내 wscan를 삭제
func (s *server) DeleteWscan(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	wscan := &wscan_struct.Wscan{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), wscan)

	wscan.Status = "deleted"
	wscan.Updated = int64(time.Now().Unix())

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), wscan)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}

	log.Printf("wscan: %v", wscan)

	// return wscan to client
	return &proto.Response{Uuid: wscan.UUID, Prompt: wscan.Prompt, Message: wscan.Message, Created: wscan.Created, Updated: wscan.Updated}, nil
}

// SendNotification 특정 유저의 wscan를 조회하여 pending 상태만 notification을 보낸다.
func (s *server) SendNotification(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "pending")
	wscans := []wscan_struct.Wscan{}
	dbClient.GetAll(context.Background(), query, &wscans)

	wg := sync.WaitGroup{}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(apns_server, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	c := apns_proto.NewAddServiceClient(conn)

	for _, wscan := range wscans {
		if wscan.NameKey == "" {
			log.Print("continue")
			continue
		}
		// request to notification grpc server
		wg.Add(1)
		go invokeNotification(c, wscan, &wg)
		// wscan의 status를 complete로 변경
		wscan.Status = "complete"
		wscan.Updated = int64(time.Now().Unix())

		_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, wscan.NameKey, nil), &wscan)
		if err != nil {
			log.Printf("Failed to put: %v", err)
		}
		log.Print("wscan notification : ", wscan.NameKey)
	}
	wg.Wait()

	return &proto.Response{}, nil
}

func (s *server) SendNotifications(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// pendding 상태의 wscan를 조회하여 notification을 보낸다.
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "pending").DistinctOn("UUID")
	wscans := []wscan_struct.Wscan{}
	dbClient.GetAll(context.Background(), query, &wscans)

	wg := sync.WaitGroup{}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(apns_server, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	c := apns_proto.NewAddServiceClient(conn)

	for _, wscan := range wscans {
		if wscan.NameKey == "" {
			log.Print("continue")
			continue
		}
		// request to notification grpc server
		wg.Add(1)
		go invokeNotification(c, wscan, &wg)
		// wscan의 status를 complete로 변경
		wscan.Status = "complete"
		wscan.Updated = int64(time.Now().Unix())

		_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, wscan.NameKey, nil), &wscan)
		if err != nil {
			log.Printf("Failed to put: %v", err)
		}
		log.Print("wscans notification : ", wscan.NameKey)
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

	// set last wscan
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
	kind := datastore.GetKindByPrefix(app+":"+env, "wscan")

	wscan := &wscan_struct.Wscan{}
	wscan.UUID = _uuid
	wscan.Prompt = prompt
	wscan.Context = gen_context
	wscan.Message = message
	wscan.Created = int64(time.Now().Unix())
	wscan.Status = "pending"
	wscan.NameKey = uuid.New().String()

	NameKey := datastore.NameKey(kind, wscan.NameKey, nil)
	_, err := dbClient.Put(context.Background(), NameKey, wscan)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}
	log.Printf("wscan: %v", wscan)
	return wscan.NameKey
}

func invokeNotification(c apns_proto.AddServiceClient, wscan wscan_struct.Wscan, wg *sync.WaitGroup) {
	ctx := context.Background()
	_, err := c.SendNotification(ctx, &apns_proto.Request{Uuid: wscan.UUID, Title: "", Subtitle: "", Body: wscan.Message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	wg.Done()
}
