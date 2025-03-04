package repositories

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/entity"
)

type FirestoreUserRepository struct {
	Users *firestore.CollectionRef
	Ctx   *context.Context
}

func NewFirestoreUserRepository(db *firestore.Client) *FirestoreUserRepository {
	ctx := context.Background()

	return &FirestoreUserRepository{
		Users: db.Collection("users"),
		Ctx:   &ctx,
	}
}

func (r *FirestoreUserRepository) Create(user *entity.User) error {
	_, err := r.Users.NewDoc().Create(
		*r.Ctx,
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *FirestoreUserRepository) FindByEmail(email string) (*entity.User, error) {
	query := r.Users.Where(
		"Email", "==", email,
	).Limit(1)

	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var user *entity.User

	err = docs[0].DataTo(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *FirestoreUserRepository) Exists(email string) (bool, error) {
	query := r.Users.Where("Email", "==", email).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return false, err
	}

	return len(docs) > 0, nil
}

func (r *FirestoreUserRepository) FindByID(id string) (*entity.User, error) {
	query := r.Users.Where("ID", "==", id).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var user *entity.User

	err = docs[0].DataTo(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *FirestoreUserRepository) Update(user *entity.User) error {
	userRef := r.Users.Doc(user.ID)

	_, err := userRef.Set(*r.Ctx, map[string]interface{}{
		"Name": user.Name,
		"Location": map[string]float64{
			"Latitude":  user.Location.Latitude,
			"Longitude": user.Location.Longitude,
		},
		"Updated_at": user.UpdatedAt,
	}, firestore.MergeAll)

	return err
}
