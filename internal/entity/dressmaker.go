package entity

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/paulozy/costurai/pkg"
)

type Address struct {
	City         string   `json:"city"`
	State        string   `json:"state"`
	Neighborhood string   `json:"neighborhood"`
	Street       string   `json:"street"`
	Number       string   `json:"number"`
	Location     Location `json:"location"`
}

type Dressmaker struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Name           string   `json:"name"`
	Contact        string   `json:"contact"`
	Enabled        bool     `json:"enabled"`
	Grade          float64  `json:"grade"`
	Services       []string `json:"services"`
	SubscriptionId *string  `json:"subscriptionId"`
	Address        Address  `json:"address"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDressmakerInput struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Contact  string   `json:"contact"`
	Services []string `json:"services"`
	Address  Address  `json:"address"`
}

type UpdateDressmakerInput struct {
	ID       string   `json:"id"`
	Name     string   `json:"name,omitempty"`
	Contact  string   `json:"contact,omitempty"`
	Address  Address  `json:"address,omitempty"`
	Services []string `json:"services,omitempty"`
}

func NewDressmaker(params CreateDressmakerInput) (*Dressmaker, error) {
	passHash, err := pkg.Encrypt(params.Password)
	if err != nil {
		return nil, err
	}

	dressmaker := &Dressmaker{
		ID:        uuid.New().String(),
		Email:     params.Email,
		Password:  string(passHash),
		Name:      params.Name,
		Contact:   params.Contact,
		Services:  params.Services,
		Grade:     0,
		Enabled:   false,
		Address:   params.Address,
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

func (dressmaker *Dressmaker) AddSubscription(sub *Subscription) {
	dressmaker.SubscriptionId = &sub.ID
}

func (dressmaker *Dressmaker) UpdateGrade(grade float64) {
	dressmaker.Grade = grade
}

func (dressmaker *Dressmaker) Update(params UpdateDressmakerInput) {
	if params.Name != "" {
		dressmaker.Name = params.Name
	}
	if params.Contact != "" {
		dressmaker.Contact = params.Contact
	}
	// Check if Address is not the zero value
	if (params.Address != Address{}) {
		dressmaker.Address = params.Address
	}
	if params.Services != nil && len(params.Services) > 0 {
		dressmaker.Services = params.Services
	}
	dressmaker.UpdatedAt = time.Now()
}

func (dressmaker *Dressmaker) CalculateGrade(reviews []Review) float64 {
	if len(reviews) == 0 {
		return 0
	}

	totalGrade := 0.0
	for _, review := range reviews {
		totalGrade += review.Grade
	}

	return math.Round(totalGrade / float64(len(reviews)))
}
