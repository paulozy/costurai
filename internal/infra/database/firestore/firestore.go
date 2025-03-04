package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(projectId string) *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Panic("Error to connect firestore", err)
	}

	return client
}
