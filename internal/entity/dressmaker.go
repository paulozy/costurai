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

	Name         string        `json:"name"`
	Contact      string        `json:"contact"`
	Enabled      bool          `json:"enabled"`
	Grade        float64       `json:"grade"`
	Services     []string      `json:"services"`
	Location     Location      `json:"location"`
	Reviews      []Review      `json:"reviews"`
	Subscription *Subscription `json:"subscription"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
		Enabled:   false,
		Reviews:   []Review{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return dressmaker, nil
}

func (dressmaker *Dressmaker) Enable() {
	dressmaker.Enabled = true
}

func (dressmaker *Dressmaker) Disable() {
	dressmaker.Enabled = false
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
	dressmaker.UpdatedAt = time.Now()
}

func (dressmaker *Dressmaker) calculateGrade() float64 {
	if len(dressmaker.Reviews) == 0 {
		return 0
	}

	totalGrade := 0.0
	for _, review := range dressmaker.Reviews {
		totalGrade += review.Grade
	}

	return math.Round(totalGrade / float64(len(dressmaker.Reviews)))
}
