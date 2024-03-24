package datastore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

var (
	projectId       = os.Getenv("PROJECT_ID")
	env             = os.Getenv("ENV")
	dataStoreClient *datastore.Client
)

func GetClient(ctx context.Context) *datastore.Client {
	if env != "live" {
		log.Printf(env, projectId)
	}
	if dataStoreClient != nil {
		return dataStoreClient
	}
	var err error
	dataStoreClient, err = datastore.NewClient(ctx, projectId)
	if err != nil {
		log.Printf("get ds client" + err.Error())
	}
	return dataStoreClient
}

func IncompleteKey(kind string, parent *datastore.Key) *datastore.Key {
	return datastore.IncompleteKey(kind, parent)
}

func NameKey(kind, name string, parent *datastore.Key) *datastore.Key {
	return datastore.NameKey(kind, name, parent)
}

func IDKey(kind string, id int64, parent *datastore.Key) *datastore.Key {
	return datastore.IDKey(kind, id, parent)
}

func NewQuery(kind string) *datastore.Query {
	return datastore.NewQuery(kind)
}

func Close() {
	if dataStoreClient != nil {
		dataStoreClient.Close()
	}
}

func GetKindByPrefix(prefix, kind string) string {
	return prefix + ":" + kind
}
