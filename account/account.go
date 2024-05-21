package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	proto "github.com/hojin-kr/go-grpc/account/proto"
	account "github.com/hojin-kr/go-grpc/account/struct"
	"github.com/hojin-kr/go-grpc/gcp/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")
var app = os.Getenv("APP")
var (
	appleBundleID   = os.Getenv("APPLE_BUNDLE_ID")
	appleTeamID     = os.Getenv("APPLE_TEAM_ID")
	appleisgnKeyB64 = os.Getenv("APPLE_SIGN_KEY_B64")
)

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

func (s *server) Init(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// Sign Up if Request uuid is empty
	accountUUID := request.GetUuid()
	token := request.GetToken()
	platform := request.GetPlatform()
	log.Printf("action/init/params/%s/%s/%s", platform, token, accountUUID)
	if token == "" || platform == "" {
		log.Printf("error/init/params/%s/%s/%s", platform, token, accountUUID)
		return &proto.Response{Uuid: ""}, nil
	}
	// 1. platform에서 조회
	kindPlatfom := datastore.GetKindByPrefix(app+":"+env, "platform")

	platforms := []account.Platform{}

	dbClient := datastore.GetClient(context.Background())
	query := datastore.NewQuery(kindPlatfom).FilterField("Token", "=", token).FilterField("Platform", "=", platform)
	dbClient.GetAll(context.Background(), query, &platforms)

	kindAccount := datastore.GetKindByPrefix(app+":"+env, "account")
	if len(platforms) > 0 {
		account := &account.Account{}
		dbClient.Get(context.Background(), datastore.NameKey(kindAccount, platforms[0].AccountID, nil), account)
		log.Printf("action/init/success/%s/%s/%s/%s", account.UUID, account.Status, platform, token)
		return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
	}

	// 2. platform이 없으면 새로 생성
	// 2-1. account 생성
	accountUUID = uuid.New().String()
	timestampStr := strconv.Itoa(int(time.Now().Unix()))

	newAccount := account.Account{
		UUID:    accountUUID,
		Status:  "active",
		Created: timestampStr,
		Updated: timestampStr,
	}

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kindAccount, accountUUID, nil), &newAccount)
	if err != nil {
		log.Printf("error/init/put/%v", err)
		return &proto.Response{}, nil
	}

	// 2-2. platform 생성
	newPlatform := account.Platform{
		AccountID: accountUUID,
		Token:     token,
		Platform:  platform,
	}
	_, err = dbClient.Put(context.Background(), datastore.IncompleteKey(kindPlatfom, nil), &newPlatform)
	if err != nil {
		log.Printf("error/init/put/%v", err)
		return &proto.Response{}, nil
	}
	// 3. account에서 조회
	account := &account.Account{}
	dbClient.Get(context.Background(), datastore.NameKey(kindAccount, accountUUID, nil), account)
	log.Printf("action/init/new/%s/%s/%s/%s", account.UUID, account.Status, platform, token)
	log.Printf("action/init/success/%s/%s/%s/%s", account.UUID, account.Status, platform, token)
	return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
}

func (s *server) Add(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	token := request.GetToken()
	platform := request.GetPlatform()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "platform")

	newPlatform := account.Platform{
		AccountID: accountUUID,
		Token:     token,
		Platform:  platform,
	}

	_, err := dbClient.Put(context.Background(), datastore.IncompleteKey(kind, nil), &newPlatform)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("newPlatform: %v", newPlatform)

	account := &account.Account{}
	dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), account)
	return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
}

func (s *server) Delete(_ context.Context, request *proto.Request) (*proto.Response, error) {
	accountUUID := request.GetUuid()
	token := request.GetToken()
	platform := request.GetPlatform()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix(app+":"+env, "platform")

	platforms := []account.Platform{}
	query := datastore.NewQuery(kind).FilterField("AccountID", "=", accountUUID)
	keys, _ := dbClient.GetAll(context.Background(), query, &platforms)

	// update account status
	kindAccount := datastore.GetKindByPrefix(app+":"+env, "account")
	account := &account.Account{}
	dbClient.Get(context.Background(), datastore.NameKey(kindAccount, platforms[0].AccountID, nil), account)
	account.Status = "inactive"
	account.Updated = strconv.Itoa(int(time.Now().Unix()))
	dbClient.Put(context.Background(), datastore.NameKey(kindAccount, platforms[0].AccountID, nil), account)

	if len(platforms) == 0 {
		log.Printf("error/delete/notfound/%s/%s/%s", platform, token, accountUUID)
		return &proto.Response{}, nil
	}

	dbClient.Delete(context.Background(), keys[0])
	if platform == "apple" {
		revokeAppleToken(token)
	}

	log.Printf("action/delete/success/%s/%s/%s", platform, token, accountUUID)
	return &proto.Response{}, nil
}

func revokeAppleToken(token string) {
	accessToken := generateToken(token)
	if accessToken == "" {
		log.Printf("error/revokeAppleToken/%s", token)
		return
	}

	url := "https://appleid.apple.com/auth/revoke"
	req, err := http.NewRequest("POST", url, strings.NewReader("client_id="+appleBundleID+"&client_secret="+createClientSecret(appleTeamID, appleBundleID)+"&token="+accessToken+"&token_type_hint=access_token"))
	if err != nil {
		log.Printf("error/revokeAppleToken/%v", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error/revokeAppleToken/%v", err)
		return
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("error/revokeAppleToken/%v", err)
		return
	}

	log.Printf("data: %v", data)

	log.Printf("action/revokeAppleToken/success/%s", token)
}

// Generate and validate tokens
func generateToken(code string) string {
	url := "https://appleid.apple.com/auth/token"
	// Create a new request
	log.Printf("revoke/apple/generate/token/%s", code)
	appleClientSecret := createClientSecret(appleTeamID, appleBundleID)
	log.Printf("revoke/apple/generate/token/%s", appleClientSecret)
	req, err := http.NewRequest("POST", url, strings.NewReader("client_id="+appleBundleID+"&client_secret="+appleClientSecret+"&code="+code+"&grant_type=authorization_code"))
	if err != nil {
		log.Printf("error/generateToken/%v", err)
		return ""
	}

	// Set the content type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new Client
	client := &http.Client{}

	// Send the request via a client
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error/generateToken/%v", err)
		return ""
	}

	// Close the response body
	defer resp.Body.Close()

	// Decode the response body
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("error/generateToken/%v", err)
		return ""
	}

	// Print the data
	log.Printf("data: %v", data)

	// Return the access token
	return data["access_token"].(string)

}

func createClientSecret(appleTeamID, appleBundleID string) string {
	// Define the signing method and create a new token
	token := jwt.New(jwt.SigningMethodES256)

	// Set some claims
	token.Claims = jwt.MapClaims{
		"iss": appleTeamID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": appleBundleID,
	}

	// Apple private key
	appleSignKey, _ := b64.StdEncoding.DecodeString(appleisgnKeyB64)
	log.Printf("revoke/apple/createClientSecret/%s", appleSignKey)
	// Parse the private key
	block, _ := pem.Decode([]byte(appleSignKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Printf("error/createClientSecret: failed to decode PEM block containing the key")
		return ""
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("error/createClientSecret: failed to parse PKCS8 private key: %v", err)
		return ""
	}

	// Cast the parsed key to the correct type (ecdsa.PrivateKey)
	ecdsaKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		log.Printf("error/createClientSecret: parsed key is not of type *ecdsa.PrivateKey")
		return ""
	}

	// Sign and get the complete encoded token as a string
	tokenStr, err := token.SignedString(ecdsaKey)
	if err != nil {
		log.Printf("error/createClientSecret: failed to sign token: %v", err)
		return ""
	}

	// Return the token
	return tokenStr
}
