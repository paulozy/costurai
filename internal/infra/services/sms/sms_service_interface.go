package services

type OTPServiceInterface interface {
	Send(to string) error
}
