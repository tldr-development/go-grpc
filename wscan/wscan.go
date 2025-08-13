package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
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
const modelName = "gemini-2.0-flash-001"

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
	if err != nil {
		panic(err)
	}
	if env != "live" {
		log.Printf("Run server")
	}

	// if not set env, app, projectID
	if env == "" {
		panic("ENV is not set")
	}
	// panic if app is not set
	if app == "" {
		panic("APP is not set")
	}
	if projectID == "" {
		panic("PROJECT_ID is not set")
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(e)
	}
}

func (s *server) Wscan(_ context.Context, request *proto.Request) (*proto.Response, error) {
	log.Printf("request: %v", request)
	uuid := uuid.New().String()
	imageBytes := request.GetImage()

	prompt := `
	*system:
	Optimized Luggage Analysis Prompt
	Analyze the given image of an open suitcase and provide a structured breakdown of its contents and estimated weight. Follow these steps:
	Identify and List Objects
	Recognize and list all visible objects inside the suitcase.
	Be as specific as possible (e.g., jeans, striped t-shirt, headphones, passport, wallet with cash, camera).
	Estimate the Weight of Each Object
	Provide an estimated weight (in g) for each object based on typical values.
	Example: Headphones: 300g, Camera: 600g.
	Consider Additional Hidden Items
	Assume that extra clothing or smaller objects may be present in unseen parts of the suitcase.
	Estimate a reasonable additional weight based on suitcase size and object distribution.
	Example: Additional clothing in unseen areas: +1000g.
	Include the Suitcases Own Weight
	Consider the weight of the suitcase based on its type:
	Carry-on: ~3000-4000g
	Medium-sized checked luggage: ~4000-5000g
	Large suitcase: ~5000-6000g
	Calculate the Total Estimated Weight
	Sum the weights of visible items, estimated hidden items, and the suitcase itself.
	Example: Total estimated weight: 10000g.

	*response format:
	json
	{
		"type": "OBJECT",
		"properties": {
		"items_weight": {
			"type": "ARRAY",
			"items": {
			"type": "OBJECT",
			"properties": {
				"name": { "type": "STRING" },
				"weight": { "type": "INTEGER" }
			},
			"required": ["name", "weight"]
			}
		},
		"suitcase_weight": {
			"type": "INTEGER"
		},
		"hidden_items_weight": {
			"type": "INTEGER"
		}
		}
	}
	`

	messages := generateByGemini(prompt, request.GetContext(), imageBytes)

	normalized := ""
	if len(messages) > 0 {
		if s, ok := ensureWscanJSON(messages[0]); ok {
			normalized = s
		} else {
			// last attempt: re-ask strictly for JSON only
			strictPrompt := prompt + "\nReturn only JSON matching exactly the schema above. No code fences. No markdown. No explanations."
			mm := generateByGemini(strictPrompt, request.GetContext(), imageBytes)
			if len(mm) > 0 {
				if s2, ok2 := ensureWscanJSON(mm[0]); ok2 {
					normalized = s2
				}
			}
		}
	}

	// fallback minimal empty structure if still empty
	if normalized == "" {
		normalized = `{"items_weight":[],"suitcase_weight":0,"hidden_items_weight":0}`
	}

	// persist as pending for downstream notification pipeline
	_ = setWscanDatastore(uuid, prompt, request.GetContext(), normalized)

	return &proto.Response{Uuid: uuid, Prompt: prompt, Message: normalized, Created: int64(time.Now().Unix()), Updated: int64(time.Now().Unix())}, nil
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
		go generateByGemini(wscan.Prompt, wscan.Context, nil)
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

func generateByGemini(prompt string, gen_context string, imageBytes []byte) []string {
	// context with timeout for robustness
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Println("gen/client/error:", err)
		return []string{}
	}
	defer client.Close()

	gm := client.GenerativeModel(modelName)
	// safety setting (block high-risk only)
	gm.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockOnlyHigh},
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockOnlyHigh},
	}
	// lower temperature for determinism
	gm.SetTemperature(0.2)

	// 멀티모달 입력 구성
	parts := []genai.Part{
		genai.Text(prompt + "\n" + gen_context),
	}
	if imageBytes != nil {
		parts = append(parts, genai.Blob{MIMEType: "image/jpeg", Data: imageBytes})
	}

	// basic retry with backoff
	var resp *genai.GenerateContentResponse
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		resp, lastErr = gm.GenerateContent(ctx, parts...)
		if lastErr == nil {
			break
		}
		// if blocked or context deadline, do not spam
		var be *genai.BlockedError
		if errors.As(lastErr, &be) {
			log.Println("gen/blocked:", be.Error())
			break
		}
		if errors.Is(lastErr, context.DeadlineExceeded) || errors.Is(lastErr, context.Canceled) {
			log.Println("gen/timeout:", lastErr)
			break
		}
		time.Sleep(time.Duration(300*(attempt+1)) * time.Millisecond)
	}
	if lastErr != nil {
		log.Println("gen/error:", lastErr)
		return []string{}
	}

	partsRet := printResponse(resp)
	// try to extract normalized JSON from candidates
	for _, p := range partsRet {
		if s, ok := ensureWscanJSON(p); ok {
			return []string{s}
		}
	}
	log.Println("generateByGemini")
	return partsRet
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

// ensureWscanJSON tries to extract a JSON object from the given text and
// validate it against the expected wscan schema. Returns normalized (minified)
// JSON string and true when valid, or "", false when invalid.
func ensureWscanJSON(text string) (string, bool) {
	raw := extractFirstJSONObject(text)
	if raw == "" {
		return "", false
	}
	// validate shape by unmarshaling into strict struct
	type itemWeight struct {
		Name   string `json:"name"`
		Weight int    `json:"weight"`
	}
	type wscanShape struct {
		ItemsWeight       []itemWeight `json:"items_weight"`
		SuitcaseWeight    int          `json:"suitcase_weight"`
		HiddenItemsWeight int          `json:"hidden_items_weight"`
	}
	var v wscanShape
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", false
	}
	// basic required checks
	if v.ItemsWeight == nil {
		v.ItemsWeight = []itemWeight{}
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// extractFirstJSONObject finds the first top-level JSON object in text.
func extractFirstJSONObject(s string) string {
	// quick path: if text already looks like a JSON object without fences
	trimmed := strings.TrimSpace(s)
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return trimmed
	}
	// remove common Markdown code fences
	fenceRe := regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")
	if m := fenceRe.FindStringSubmatch(s); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	// fallback: scan for first '{' and find matching closing '}'
	i := strings.Index(s, "{")
	if i == -1 {
		return ""
	}
	depth := 0
	inStr := false
	esc := false
	for j := i; j < len(s); j++ {
		c := s[j]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		if c == '"' {
			inStr = true
			continue
		}
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				return strings.TrimSpace(s[i : j+1])
			}
		}
	}
	return ""
}

func setWscanDatastore(_uuid, prompt, gen_context, message string) string {
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
