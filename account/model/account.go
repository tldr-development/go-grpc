// account/model/account.go
package model

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"account/proto"

	"gcp/datastore"

	"github.com/google/uuid"
)

type Account struct {
	UUID    string // uuid
	Status  string // status
	Created string // created at timestamp
	Updated string // updated at timestamp
}

type Platform struct {
	AccountID string // Account uuid id (datastore id)
	Token     string // platform token
	Platform  string // platform name (ex. github, google, kakao)
}

type Server struct {
	proto.UnimplementedAddServiceServer
}

var env = os.Getenv("ENV")

func (s *Server) Init(_ context.Context, request *proto.Request) (*proto.Response, error) {
	// Sign Up if Request uuid is empty
	accountUUID := request.GetUuid()

	dbClient := datastore.GetClient(context.Background())
	kind := datastore.GetKindByPrefix("test", "account")

	if accountUUID != "" {
		log.Printf("account_uuid: %s", accountUUID)
		// DB에서 uuid로 조회
		account := &Account{}
		dbClient.Get(context.Background(), datastore.NameKey(kind, accountUUID, nil), account)
		return &proto.Response{Uuid: account.UUID, Status: account.Status, Created: account.Created, Updated: account.Updated}, nil
	}
	// Sign Up
	accountUUID = uuid.New().String()
	log.Printf("account_uuid: %s", accountUUID)

	timestampStr := strconv.Itoa(int(time.Now().Unix()))

	newAccount := Account{
		UUID:    accountUUID,
		Status:  "active",
		Created: timestampStr,
		Updated: timestampStr,
	}

	_, err := dbClient.Put(context.Background(), datastore.NameKey(kind, accountUUID, nil), &newAccount)
	if err != nil {
		log.Printf("Failed to put: %v", err)
		return &proto.Response{}, nil
	}

	log.Printf("newAccount: %v", newAccount)
	return &proto.Response{Uuid: newAccount.UUID, Status: newAccount.Status, Created: newAccount.Created, Updated: newAccount.Updated}, nil
}
