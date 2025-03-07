package services

import "time"

type OTPServiceInterface interface {
	SendVerification(phone string) (*OTPServiceSendVerificationOutput, error)
	Verify(code, phone string) (*OTPServiceVerifyOutput, error)
}

type SendCodeAttempt struct {
	Time       time.Time `json:"time"`
	Channel    string    `json:"channel"`
	AttemptSID string    `json:"attempt_sid"`
}

type OTPServiceSendVerificationOutput struct {
	SID          string            `json:"sid"`
	ServiceSID   string            `json:"service_sid"`
	AccountSID   string            `json:"account_sid"`
	To           string            `json:"to"`
	Channel      string            `json:"channel"`
	Status       string            `json:"status"`
	Valid        bool              `json:"valid"`
	DateCreated  time.Time         `json:"date_created"`
	DateUpdated  time.Time         `json:"date_updated"`
	Lookup       map[string]string `json:"lookup"`
	Amount       *string           `json:"amount"`
	Payee        *string           `json:"payee"`
	SendAttempts []SendCodeAttempt `json:"send_code_attempts"`
	SNA          *string           `json:"sna"`
	URL          string            `json:"url"`
}

type OTPServiceVerifyOutput struct {
	SID                   string    `json:"sid"`
	ServiceSID            string    `json:"service_sid"`
	AccountSID            string    `json:"account_sid"`
	To                    string    `json:"to"`
	Channel               string    `json:"channel"`
	Status                string    `json:"status"`
	Valid                 bool      `json:"valid"`
	Amount                *string   `json:"amount"`
	Payee                 *string   `json:"payee"`
	SNAAttemptsErrorCodes []string  `json:"sna_attempts_error_codes"`
	DateCreated           time.Time `json:"date_created"`
	DateUpdated           time.Time `json:"date_updated"`
}
