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
	doc, err := r.Subscriptions.Doc(id).Get(*r.Ctx)
	if err != nil {
		return nil, err
	}
	var sub entity.Subscription
	if err := doc.DataTo(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *FirestoreSubscriptionRepository) Update(subscription *entity.Subscription) error {
	subscription.UpdatedAt = time.Now()
	_, err := r.Subscriptions.Doc(subscription.ID).Set(*r.Ctx, subscription)
	return err
}
