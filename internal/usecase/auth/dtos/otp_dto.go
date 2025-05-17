package dtos

type OTPSendAndVerifyInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code,omitempty"`
}
