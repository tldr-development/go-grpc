package cloudstorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"encoding/base64"

	"cloud.google.com/go/storage"
)

// google cloud platform cloud storage client

var (
	googleAccessID = os.Getenv("G_ACCESS_ID")
	b64PrivateKey  = os.Getenv("G_PRIVATE_KEY_B64")
)

func GenerateV4GetObjectSignedURL(w io.Writer, bucket, object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()
	privateKey, _ := base64.StdEncoding.DecodeString(b64PrivateKey)
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: googleAccessID,
		PrivateKey:     []byte(privateKey),
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	return u, nil
}
