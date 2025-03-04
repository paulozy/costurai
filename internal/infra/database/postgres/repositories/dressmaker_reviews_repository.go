package repositories

import (
	"database/sql"

	"github.com/paulozy/costurai/internal/entity"
)

type DressmakerReviewsRepository struct {
	DB *sql.DB
}

func NewDressmakerReviewsRepository(db *sql.DB) *DressmakerReviewsRepository {
	return &DressmakerReviewsRepository{
		DB: db,
	}
}

func (r *DressmakerReviewsRepository) Create(review *entity.Review) error {
	stmt, err := r.DB.Prepare(
		`
			INSERT INTO 
				dressmakers_reviews (
					id, 
					dressmaker_id, 
					user_id, 
					comment,
					grade,
					created_at, 
					updated_at
				) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(review.ID, review.DressmakerID, review.UserID, review.Comment, review.Grade, review.CreatedAt, review.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
