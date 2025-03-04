package services

import (
	"fmt"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioService struct {
	Client  twilio.RestClient
	SID     string
	Channel string
}

const channel string = "sms"

func NewTwilioService(sid string) *TwilioService {
	client := twilio.NewRestClient()

	return &TwilioService{
		Client:  *client,
		SID:     sid,
		Channel: channel,
	}
}

func (ts *TwilioService) SendVerification(phone string) (*verify.VerifyV2Verification, error) {
	params := &verify.CreateVerificationParams{
		To:      &phone,
		Channel: &ts.Channel,
	}

	res, err := ts.Client.VerifyV2.CreateVerification(
		ts.SID,
		params,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return res, nil
}

func (ts *TwilioService) Verify(code string) (*verify.VerifyV2VerificationCheck, error) {
	params := &verify.CreateVerificationCheckParams{
		Code: &code,
	}

	res, err := ts.Client.VerifyV2.CreateVerificationCheck(
		ts.SID,
		params,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return res, nil
}
