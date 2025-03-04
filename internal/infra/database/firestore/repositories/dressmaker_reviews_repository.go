package repositories

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/entity"
)

type DressmakerReviewsRepository struct {
	Reviews *firestore.CollectionRef
	Ctx     *context.Context
}

func NewFirestoreReviewsRepository(db *firestore.Client) *DressmakerReviewsRepository {
	ctx := context.Background()

	return &DressmakerReviewsRepository{
		Reviews: db.Collection("reviews"),
		Ctx:     &ctx,
	}
}

func (r *DressmakerReviewsRepository) Create(review *entity.Review) error {
	_, err := r.Reviews.NewDoc().Create(*r.Ctx, review)
	if err != nil {
		return err
	}

	return nil
}
