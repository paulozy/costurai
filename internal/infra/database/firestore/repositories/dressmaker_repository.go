package repositories

import (
	"context"
	"fmt"
	"math"

	"cloud.google.com/go/firestore"
	"github.com/paulozy/costurai/internal/entity"
	"github.com/paulozy/costurai/pkg"
)

type FirestoreDressmakerRepository struct {
	Dressmakers *firestore.CollectionRef
	Ctx         *context.Context
}

func NewFirestoreDressmakerRepository(db *firestore.Client) *FirestoreDressmakerRepository {
	ctx := context.Background()

	return &FirestoreDressmakerRepository{
		Dressmakers: db.Collection("dressmakers"),
		Ctx:         &ctx,
	}
}

func (r *FirestoreDressmakerRepository) Create(dressmaker *entity.Dressmaker) error {
	_, err := r.Dressmakers.NewDoc().Create(*r.Ctx, dressmaker)
	return err
}

func (r *FirestoreDressmakerRepository) FindByEmail(email string) (*entity.Dressmaker, error) {
	query := r.Dressmakers.Where(
		"Email", "==", email,
	).Limit(1)

	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var dressmaker *entity.Dressmaker

	err = docs[0].DataTo(&dressmaker)
	if err != nil {
		return nil, err
	}

	return dressmaker, nil
}

func (r *FirestoreDressmakerRepository) Exists(email string) (bool, error) {
	query := r.Dressmakers.Where("Email", "==", email).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return false, err
	}

	return len(docs) > 0, nil
}

func (r *FirestoreDressmakerRepository) FindByID(id string) (*entity.Dressmaker, error) {
	query := r.Dressmakers.Where("ID", "==", id).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, nil
	}

	var dressmaker *entity.Dressmaker

	err = docs[0].DataTo(&dressmaker)
	if err != nil {
		return nil, err
	}

	return dressmaker, nil
}

func (r *FirestoreDressmakerRepository) FindByProximity(latitude, longitude float64, maxDistance int) ([]entity.Dressmaker, error) {
	fmt.Println("maxDistance", maxDistance)

	minLat, maxLat, minLng, maxLng := getBoundingBox(latitude, longitude, float64(maxDistance))

	query := r.Dressmakers.
		Where("Enabled", "==", true).
		Where("Address.Location.Latitude", ">=", minLat).
		Where("Address.Location.Latitude", "<=", maxLat).
		Where("Address.Location.Longitude", ">=", minLng).
		Where("Address.Location.Longitude", "<=", maxLng)

	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	dressmakers := make([]entity.Dressmaker, 0, 50)
	for _, doc := range docs {
		var dressmaker *entity.Dressmaker
		err = doc.DataTo(&dressmaker)
		if err != nil {
			return nil, err
		}

		dist := pkg.HaversineDistance(
			latitude,
			longitude,
			dressmaker.Address.Location.Latitude,
			dressmaker.Address.Location.Longitude,
		)

		if dist <= float64(maxDistance) {
			dressmakers = append(dressmakers, *dressmaker)
		}
	}

	return dressmakers, nil
}

func (r *FirestoreDressmakerRepository) Update(dressmaker *entity.Dressmaker) error {
	query := r.Dressmakers.Where("ID", "==", dressmaker.ID).Limit(1)
	docs, err := query.Documents(*r.Ctx).GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		return fmt.Errorf("no dressmaker found with ID: %s", dressmaker.ID)
	}

	dressmakerRef := docs[0].Ref

	_, err = dressmakerRef.Set(*r.Ctx, map[string]interface{}{
		"Name":    dressmaker.Name,
		"Email":   dressmaker.Email,
		"Contact": dressmaker.Contact,
		"Address": map[string]interface{}{
			"Street":  dressmaker.Address.Street,
			"Number":  dressmaker.Address.Number,
			"City":    dressmaker.Address.City,
			"State":   dressmaker.Address.State,
			"Zipcode": dressmaker.Address.Zipcode,
			"Location": map[string]float64{
				"Latitude":  dressmaker.Address.Location.Latitude,
				"Longitude": dressmaker.Address.Location.Longitude,
			},
		},
		"Services":       dressmaker.Services,
		"Grade":          dressmaker.Grade,
		"SubscriptionId": dressmaker.SubscriptionId,
		"Enabled":        dressmaker.Enabled,
		"CreatedAt":      dressmaker.CreatedAt,
		"UpdatedAt":      dressmaker.UpdatedAt,
	}, firestore.MergeAll)

	return err
}

func getBoundingBox(lat, lng, radiusMeters float64) (float64, float64, float64, float64) {
	const earthRadiusMeters = 6371000.0 // Earth's radius in meters

	latDelta := (radiusMeters / earthRadiusMeters) * (180 / math.Pi)
	lngDelta := (radiusMeters / earthRadiusMeters) * (180 / math.Pi) / math.Cos(lat*math.Pi/180)

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLng := lng - lngDelta
	maxLng := lng + lngDelta

	return minLat, maxLat, minLng, maxLng
}
