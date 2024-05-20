package entity

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/paulozy/costurai/pkg"
)

type Dressmaker struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Name      string   `json:"name"`
	Contact   string   `json:"contact"`
	Location  Location `json:"location"`
	Services  []string `json:"services"`
	Grade     float64  `json:"grade"`
	Reviews   []Review `json:"reviews"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func NewDressmaker(email, password, name, contact string, location Location, services []string) (*Dressmaker, error) {
	passHash, err := pkg.Encrypt(password)
	if err != nil {
		return nil, err
	}

	dressmaker := &Dressmaker{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(passHash),
		Name:      name,
		Contact:   contact,
		Location:  location,
		Services:  services,
		Grade:     0,
		Reviews:   []Review{},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return dressmaker, nil
}

func (dressmaker *Dressmaker) AddReview(review Review) {
	dressmaker.Reviews = append(dressmaker.Reviews, review)
	dressmaker.UpdateGrade(dressmaker.calculateGrade())
}

func (dressmaker *Dressmaker) UpdateGrade(grade float64) {
	dressmaker.Grade = grade
}

func (dressmaker *Dressmaker) UpdateLocation(location Location) {
	dressmaker.Location = location
}

func (dressmaker *Dressmaker) UpdateDressmaker(name, contact string, location Location, services []string) {
	dressmaker.Name = name
	dressmaker.Contact = contact
	dressmaker.Location = location
	dressmaker.Services = services
	dressmaker.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (dressmaker *Dressmaker) calculateGrade() float64 {
	if len(dressmaker.Reviews) == 0 {
		return 0
	}

	totalGrade := 0.0
	for _, review := range dressmaker.Reviews {
		totalGrade += review.Grade
	}

	return math.Round(totalGrade/float64(len(dressmaker.Reviews)))
}
