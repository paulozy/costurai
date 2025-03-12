package repositories

import (
	"context"
	"log"

	gfirestore "cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/database/firestore"
)

type FirestoreDressmakerReviewsRepository struct {
	Reviews *gfirestore.CollectionRef
	Ctx     *context.Context
}

func NewFirestoreReviewsRepository(db database.DatabaseInterface) *FirestoreDressmakerReviewsRepository {
	ctx := context.Background()

	firestoreDB, ok := db.(*firestore.FirestoreDatabase)
	if !ok {
		log.Panic("Expected *FirestoreDatabase, got different type")
	}

	return &FirestoreDressmakerReviewsRepository{
		Reviews: firestoreDB.Client.Collection("reviews"),
		Ctx:     &ctx,
	}
}

func (r *FirestoreDressmakerReviewsRepository) Create(review *entity.Review) error {
	_, err := r.Reviews.NewDoc().Create(*r.Ctx, review)
	if err != nil {
		return err
	}

	return nil
}
