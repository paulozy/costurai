package entity

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID						string  `json:"id"`
	DressmakerID 	string  `json:"dressmakerId"`
	UserID      	string  `json:"userId"`
	Grade       	float64 `json:"grade"`
	Comment 			string  `json:"comment"`
	CreatedAt   	string  `json:"created_at"`
	UpdatedAt   	string  `json:"updated_at"`
}

func NewReview(dressmakerID, userID, comment string, grade float64) *Review {
	return &Review{
		ID:          		uuid.New().String(),
		DressmakerID: 	dressmakerID,
		UserID:      		userID,
		Grade:       		grade,
		Comment: 				comment,
		CreatedAt:   		time.Now().Format(time.RFC3339),
		UpdatedAt:   		time.Now().Format(time.RFC3339),
	}
}