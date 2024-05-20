package database

import "github.com/paulozy/costurai/internal/entity"

type DressmakerRepositoryInterface interface {
	Create(dressmaker *entity.Dressmaker) error
	FindByEmail(email string) (*entity.Dressmaker, error)
	Exists(email string) (bool, error)
	FindByID(id string) (*entity.Dressmaker, error)
	FindByProximity(latitude, longitude float64, maxDistance int) ([]entity.Dressmaker, error)
	FindByServices(services string) ([]entity.Dressmaker, error)
	Update(dressmaker *entity.Dressmaker) error
}

type UserRepositoryInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	Exists(email string) (bool, error)
	FindByID(id string) (*entity.User, error)
	Update(user *entity.User) error
}

type DressmakerReviewsRepositoryInterface interface {
	Create(review *entity.Review) error
}