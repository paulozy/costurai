package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type FirestoreDatabase struct {
	Client *firestore.Client
}

func NewFirestoreClient(projectId string) *FirestoreDatabase {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Panic("Error to connect firestore", err)
	}

	return &FirestoreDatabase{
		Client: client,
	}
}

func (f *FirestoreDatabase) Ping(ctx context.Context) error {
	_, _, err := f.Client.Collection("healthcheck").Add(ctx, map[string]interface{}{
		"status": "ok",
	})
	return err
}

func (f *FirestoreDatabase) Close() error {
	return f.Client.Close()
}
