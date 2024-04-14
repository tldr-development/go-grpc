package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/google/uuid"
	apns_proto "github.com/hojin-kr/fiber-grpc/apns/proto"
	"github.com/hojin-kr/fiber-grpc/gcp/datastore"
	proto "github.com/hojin-kr/fiber-grpc/inspire/proto"
	inspire_struct "github.com/hojin-kr/fiber-grpc/inspire/struct"
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
	_uuid := request.GetUuid()

	// 유저 정보에서 유저가 얼마 주기로 inspire를 받아야 하는지 조회
	info := getInspireInfo(_uuid)

	// info가 없으면 초기화
	if info.UUID == "" {
		info = initInspireInfo(_uuid)
		updateInspireInfo(info.UUID, info.Status, info.NotiPeriod, info.MessageLengthType, info.MessageType, info.Context, info.UserContext, info.LastMessage, info.Updated)
	}

	if info.Status == "inactive" {
		return &proto.Response{}, nil
	}
	if info.Status == "blocked" {
		return &proto.Response{}, nil
	}

	inspire := inspire_struct.Inspire{}
	gen_context := getContext(inspire.UUID)
	// inspire 생성
	generateByGemini(prompt, gen_context, _uuid)

	// invoke notification
	c := apns_proto.NewAddServiceClient(nil)
	wg := sync.WaitGroup{}

	go invokeNotification(c, inspire, &wg)

	return &proto.Response{}, nil
}

func initInspireInfo(uuid string) inspire_struct.Info {
	return inspire_struct.Info{
		UUID:              uuid,
		Status:            "active",
		NotiPeriod:        "0",
		MessageLengthType: "short",
		MessageType:       "counselor",
		Context:           "",
		UserContext:       "",
		LastMessage:       "",
		Updated:           int64(time.Now().Unix()),
		Language:          "korean",
		NameKey:           "",
	}
}

// 내 inspire 목록을 조회
func (s *server) GetInspires(_ context.Context, request *proto.Request) (*proto.Responses, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "inspire")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "complete")
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
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "inspire_last")

	inspire := &inspire_struct.Inspire{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), inspire)
	response := &proto.Response{Uuid: inspire.UUID, Prompt: inspire.Prompt, Message: inspire.Message, Created: inspire.Created, Updated: inspire.Updated}
	log.Printf("inspire: %v", inspire)

	// return inspire to client
	return response, nil
}

// 특정 시간 이후의 inspire last를 조회해서 inspire를 생성한다.
func (s *server) GenerateInspireAfterCreatedLast(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "inspire_last")

	query := datastore.NewQuery(kind).FilterField("Created", "<", request.GetCreated()).Limit(10000)
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	for _, inspire := range inspires {
		// 유저 정보에서 유저가 얼마 주기로 inspire를 받아야 하는지 조회
		info := getInspireInfo(inspire.UUID)

		// info가 없으면 초기화
		if info.UUID == "" {
			info = initInspireInfo(inspire.UUID)
			go updateInspireInfo(info.UUID, info.Status, info.NotiPeriod, info.MessageLengthType, info.MessageType, info.Context, info.UserContext, info.LastMessage, info.Updated)
		}

		if info.Status == "inactive" {
			continue
		}
		if info.Status == "blocked" {
			continue
		}
		// info에서 주기를 확인, period is second
		period, _ := strconv.Atoi(info.NotiPeriod)
		created := inspire.Created + int64(period)
		if created > int64(time.Now().Unix()) {
			continue
		}

		inspire.Context = getContext(inspire.UUID)
		generateByGemini(inspire.Prompt, inspire.Context, inspire.UUID)
	}

	return &proto.Response{}, nil
}

func getContext(_uuid string) string {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "info")

	info := &inspire_struct.Info{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, _uuid, nil), info)

	info.Context = "User Custom Context : " + info.UserContext + "\n" +
		"Context : " + info.Context + "\n" +
		"Last Message : " + info.LastMessage + "\n" +
		"Language : " + info.Language + "\n" +
		"Message Type : " + info.MessageType + "\n" +
		"Message Length Type : " + info.MessageLengthType + "\n"

	return info.Context
}

// 내 inspire를 삭제
func (s *server) DeleteInspire(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "inspire")

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
	kind := datastore.GetKindByPrefix(getKind(), "inspire")

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
	kind := datastore.GetKindByPrefix(getKind(), "inspire")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "pending")
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

// GetInspireInfo 특정 유저의 inspire 정보를 조회한다.
func (s *server) GetInspireInfo(_ context.Context, request *proto.RequestInfo) (*proto.ResponseInfo, error) {
	info := getInspireInfo(request.GetUuid())

	// return inspire info to client
	return &proto.ResponseInfo{Uuid: info.UUID, Status: info.Status, NotiPeriod: info.NotiPeriod, MessageLengthType: info.MessageLengthType, MessageType: info.MessageType, Context: info.Context, UserContext: info.UserContext, LastMessage: info.LastMessage, Updated: info.Updated}, nil
}

// getinspireinfo 특정 유저의 inspire 정보를 조회한다.
func getInspireInfo(uuid string) inspire_struct.Info {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "info")

	info := &inspire_struct.Info{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, uuid, nil), info)

	return *info
}

// UpdateInspireInfo 특정 유저의 inspire 정보를 업데이트한다.
func (s *server) UpdateInspireInfo(_ context.Context, request *proto.RequestInfo) (*proto.ResponseInfo, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "info")

	info := &inspire_struct.Info{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), info)

	info.Status = request.GetStatus()
	info.NotiPeriod = request.GetNotiPeriod()
	info.MessageLengthType = request.GetMessageLengthType()
	info.MessageType = request.GetMessageType()
	info.Context = request.GetContext()
	info.UserContext = request.GetUserContext()
	info.LastMessage = request.GetLastMessage()
	info.Updated = int64(time.Now().Unix())

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, request.GetUuid(), nil), info)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}

	log.Printf("info: %v", info)

	// return inspire info to client
	return &proto.ResponseInfo{Uuid: info.UUID, Status: info.Status, NotiPeriod: info.NotiPeriod, MessageLengthType: info.MessageLengthType, MessageType: info.MessageType, Context: info.Context, UserContext: info.UserContext, LastMessage: info.LastMessage, Updated: info.Updated}, nil
}

func updateInspireInfo(uuid string, status string, noti_period string, message_length_type string, message_type string, _context string, user_context string, last_message string, updated int64) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "info")

	info := &inspire_struct.Info{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, uuid, nil), info)

	info.Status = status
	info.NotiPeriod = noti_period
	info.MessageLengthType = message_length_type
	info.MessageType = message_type
	info.Context = _context
	info.UserContext = user_context
	info.LastMessage = last_message
	info.Updated = updated

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, uuid, nil), info)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}

	log.Printf("info: %v", info)
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
	inspireLastNameKey := ""
	for _, part := range parts {
		fmt.Println(part + "\n")
		inspireLastNameKey = setInpireDatastore(_uuid, prompt, gen_context, part)
	}
	if inspireLastNameKey != "" {
		// set last inspire
		setLastInspire(_uuid, prompt, gen_context, parts[len(parts)-1])

		// set inspire info
		info := getInspireInfo(_uuid)
		info.LastMessage = parts[len(parts)-1]
		info.Updated = int64(time.Now().Unix())
		go updateInspireInfo(info.UUID, info.Status, info.NotiPeriod, info.MessageLengthType, info.MessageType, info.Context, info.UserContext, info.LastMessage, info.Updated)
	}
	return parts
}

// set last inspire
func setLastInspire(_uuid, prompt, gen_context, message string) string {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(getKind(), "inspire_last")

	inspire := &inspire_struct.Inspire{}
	inspire.UUID = _uuid
	inspire.Prompt = prompt
	inspire.Context = gen_context
	inspire.Message = message
	inspire.Created = int64(time.Now().Unix())
	inspire.Status = "pending"
	inspire.NameKey = uuid.New().String()

	UUID := datastore.NameKey(kind, inspire.UUID, nil)
	_, err := dbClient.Put(context.Background(), UUID, inspire)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}
	log.Printf("inspire: %v", inspire)
	return inspire.NameKey
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
	kind := datastore.GetKindByPrefix(getKind(), "inspire")

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
	_, err := c.SendNotification(ctx, &apns_proto.Request{Uuid: inspire.UUID, Title: "", Subtitle: inspire.Prompt, Body: inspire.Message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	wg.Done()
}

func getKind() string {
	return app + ":" + env
}
