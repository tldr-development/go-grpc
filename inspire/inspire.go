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
	wallet_proto "github.com/hojin-kr/go-grpc/wallet/proto"
	wallet_struct "github.com/hojin-kr/go-grpc/wallet/struct"
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
var wallet_server = os.Getenv("WALLET_SERVER")

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

	// wallet ticket 감소
	if !decrWalletTicket(_uuid) {
		log.Println("decrWalletTicket failed")
		return &proto.Response{}, nil
	}

	messages := generateByGemini(prompt, gen_context)

	if len(messages) > 0 {
		for _, message := range messages {
			setInpireDatastore(_uuid, prompt, gen_context, message, "complete", nil)
		}
	}
	log.Println("Inspire")

	return &proto.Response{Uuid: _uuid, Prompt: prompt, Message: messages[0], Created: int64(time.Now().Unix()), Updated: int64(time.Now().Unix())}, nil
}

// 내 inspire 목록을 조회
func (s *server) GetInspires(_ context.Context, request *proto.Request) (*proto.Responses, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("UUID", "=", request.GetUuid()).FilterField("Status", "=", "complete").Order("-Created").Limit(10000)
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

// 마지막 updated이후로 3일 이내에 생성된 inspire를 조회하여 대상자로 선정
func (s *server) GenerateInspireAfterCreatedLast(_ context.Context, request *proto.Request) (*proto.Response, error) {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	query := datastore.NewQuery(kind).FilterField("Status", "=", "complete").FilterField("Updated", ">", request.GetCreated()).Order("-Updated").FilterField("Context", "!=", "auto")
	inspires := []inspire_struct.Inspire{}
	dbClient.GetAll(context.Background(), query, &inspires)

	// distinct inspire by UUID
	inspireMap := make(map[string]inspire_struct.Inspire)

	for _, inspire := range inspires {
		if inspireMap[inspire.UUID].Created < inspire.Created {
			inspireMap[inspire.UUID] = inspire
		}
	}
	// user count log
	log.Println("GenerateInspireAfterCreatedLast,inspireMap,", len(inspireMap))

	type MessageByPrompt struct {
		Prompt  string
		Context string
		Message string
	}

	promptMessageMap := make(map[string]MessageByPrompt)

	for _, inspire := range inspireMap {
		if promptMessageMap[inspire.Prompt].Message == "" {
			messages := generateByGemini(inspire.Prompt, inspire.Context)
			log.Println(inspire.Prompt, request.Context, messages)
			// check if message is empty
			if len(messages) == 0 {
				continue
			}
			for _, message := range messages {
				promptMessageMap[inspire.Prompt] = MessageByPrompt{
					Prompt:  inspire.Prompt,
					Context: request.Context,
					Message: message,
				}
			}
		}
	}
	// prompt count log
	log.Println("GenerateInspireAfterCreatedLast,promptMessageMap,", len(promptMessageMap))

	wg := sync.WaitGroup{}

	// set inspire to datastore
	for _, inspire := range inspireMap {
		if promptMessageMap[inspire.Prompt].Message == "" {
			continue
		}
		wg.Add(1)
		go setInpireDatastore(inspire.UUID, inspire.Prompt, "auto", promptMessageMap[inspire.Prompt].Message, "pending", &wg)
	}
	wg.Wait()

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
		if inspire.Context == "auto" {
			inspire.Updated = 0
		} else {
			inspire.Updated = int64(time.Now().Unix())
		}

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

	query := datastore.NewQuery(kind).FilterField("Status", "=", "pending")
	inspires := []inspire_struct.Inspire{}

	dbClient.GetAll(context.Background(), query, &inspires)

	// distinct inspire by UUID
	inspireMap := make(map[string]inspire_struct.Inspire)
	for _, inspire := range inspires {
		inspireMap[inspire.UUID] = inspire
	}

	wg := sync.WaitGroup{}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(apns_server, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	c := apns_proto.NewAddServiceClient(conn)

	for _, inspire := range inspireMap {
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

func generateByGemini(prompt string, gen_context string) []string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(model)
	// safty setting
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}
	model.SetTemperature(0.9)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt+"\n"+gen_context))
	if err != nil {
		log.Println("gen/error/" + prompt + "\n" + gen_context)
		log.Println(err)
	}

	parts := printResponse(resp)
	log.Println("generateByGemini")
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

func setInpireDatastore(_uuid, prompt, gen_context, message, status string, wg *sync.WaitGroup) string {
	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "inspire")

	inspire := &inspire_struct.Inspire{}
	inspire.UUID = _uuid
	inspire.Prompt = prompt
	inspire.Context = gen_context
	inspire.Message = message
	inspire.Created = int64(time.Now().Unix())
	inspire.Updated = int64(time.Now().Unix())
	inspire.Status = status
	inspire.NameKey = uuid.New().String()

	NameKey := datastore.NameKey(kind, inspire.NameKey, nil)
	_, err := dbClient.Put(context.Background(), NameKey, inspire)
	if err != nil {
		log.Printf("Failed to put: %v", err)
	}
	log.Printf("setInpireDatastore: %v", inspire)
	if wg != nil {
		wg.Done()
	}
	return inspire.NameKey
}

func invokeNotification(c apns_proto.AddServiceClient, inspire inspire_struct.Inspire, wg *sync.WaitGroup) {
	ctx := context.Background()
	_, err := c.SendNotification(ctx, &apns_proto.Request{Uuid: inspire.UUID, Title: "", Subtitle: "", Body: inspire.Message})
	if err != nil {
		log.Println("error/invokeNoti/%w", err)
	}
	wg.Done()
}

// wallet의 ticket을 조회하고 감소
func decrWalletTicket(_uuid string) bool {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(wallet_server, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Println("error/inspire/decr/wallet/ticket/%w", err)
	}

	c := wallet_proto.NewAddServiceClient(conn)
	req := &wallet_proto.Request{Uuid: _uuid}
	res, err := c.Get(context.Background(), req)
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}

	_wallet := &wallet_struct.Wallet{}
	_wallet.UUID = res.Uuid
	_wallet.Ticket = res.Ticket

	if _wallet.UUID == "" {
		log.Println("wallet not found uuid: ", _uuid)
		return false
	}

	if _wallet.Ticket == 0 {
		log.Println("wallet ticket is 0 uuid: ", _uuid)
		return false
	}

	_wallet.Ticket = _wallet.Ticket - 1
	_wallet.Updated = int64(time.Now().Unix())

	_, err = c.Update(context.Background(), &wallet_proto.Request{Uuid: _wallet.UUID, Ticket: _wallet.Ticket, Status: "active"})
	if err != nil {
		return false
	}

	log.Printf("decrWalletTicket uuid: %v, ticket: %v", _uuid, _wallet.Ticket)
	return true
}
