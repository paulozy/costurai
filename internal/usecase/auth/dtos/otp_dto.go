package dtos

type SendOTPInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code,omitempty"`
}

type VerifyOTPInput struct {
	Phone        string `json:"phone"`
	Code         string `json:"code"`
	Enabling     string `json:"enabling"`
	UserID       string `json:"userId,omitempty"`
	DressmakerID string `json:"dressmakerId,omitempty"`
}
