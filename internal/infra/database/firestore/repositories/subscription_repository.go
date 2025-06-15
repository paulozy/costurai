package repositories

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/entity"
)

type FirestoreSubscriptionRepository struct {
	Subscriptions *firestore.CollectionRef
	Ctx           *context.Context
}

func NewFirestoreSubscriptionRepository(db *firestore.Client) *FirestoreSubscriptionRepository {
	ctx := context.Background()

	return &FirestoreSubscriptionRepository{
		Subscriptions: db.Collection("subscriptions"),
		Ctx:           &ctx,
	}
}

func (r *FirestoreSubscriptionRepository) Create(subscription *entity.Subscription) error {
	_, err := r.Subscriptions.NewDoc().Create(*r.Ctx, subscription)

	if err != nil {
		return err
	}

	return nil
}

func (r *FirestoreSubscriptionRepository) FindByID(id string) (*entity.Subscription, error) {
	query := r.Subscriptions.Where("ID", "==", id).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var subscription entity.Subscription
	err = docs[0].DataTo(&subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *FirestoreSubscriptionRepository) Update(subscription *entity.Subscription) error {
	subscription.UpdatedAt = time.Now()
	_, err := r.Subscriptions.Doc(subscription.ID).Set(*r.Ctx, subscription)
	return err
}
