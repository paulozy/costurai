package services

type OTPServiceInterface interface {
	Send(to string) error
	Verify(phone string, code string) (bool, error)
}
