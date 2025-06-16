package services

import (
	"fmt"
	"time"

	"github.com/paulozy/costurai/configs"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioService struct {
	Client     *twilio.RestClient
	ServiceSID string
	Channel    string
}

func NewTwilioService() *TwilioService {
	configs, _ := configs.LoadConfig(configs.Env)

	params := twilio.ClientParams{
		Username: configs.TwilioSID,
		Password: configs.TwilioAuthToken,
	}

	client := twilio.NewRestClientWithParams(params)
	client.SetTimeout(time.Duration(configs.SMSTimeout) * time.Second)

	return &TwilioService{
		Client:     client,
		ServiceSID: configs.TwilioSMSServiceSID,
		Channel:    configs.TwilioChannel,
	}
}

func (s *TwilioService) Send(to string) error {
	params := &verify.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	_, err := s.Client.VerifyV2.CreateVerification(s.ServiceSID, params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *TwilioService) Verify(phone, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)

	resp, err := s.Client.VerifyV2.CreateVerificationCheck(s.ServiceSID, params)
	if err != nil {
		return false, err
	} else if resp.Status == nil || *resp.Status != "approved" {
		return false, err
	}

	return true, nil
}
