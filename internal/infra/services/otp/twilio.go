package services

import (
	"fmt"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioService struct {
	Client    twilio.RestClient
	ServiceID string
	Channel   string
}

const channel string = "sms"

func NewTwilioService(accountSid, authToken, serviceId string) *TwilioService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid, // Seu Account SID
		Password: authToken,  // Seu Auth Token
	})

	return &TwilioService{
		Client:    *client,
		ServiceID: serviceId,
		Channel:   channel,
	}
}

func (ts *TwilioService) SendVerification(phone string) (*OTPServiceSendVerificationOutput, error) {
	params := &verify.CreateVerificationParams{
		To:      &phone,
		Channel: &ts.Channel,
	}

	res, err := ts.Client.VerifyV2.CreateVerification(
		ts.ServiceID,
		params,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	output := &OTPServiceSendVerificationOutput{
		SID:        *res.Sid,
		ServiceSID: *res.ServiceSid,
		AccountSID: *res.AccountSid,
		To:         *res.To,
		Channel:    *res.Channel,
		Status:     *res.Status,
		Valid:      *res.Valid,
		URL:        *res.Url,
	}

	return output, nil
}

func (ts *TwilioService) Verify(code, phone string) (*OTPServiceVerifyOutput, error) {
	fmt.Println("code: %s, phone: %s", code, phone)

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)

	fmt.Println(*params.Code, *params.To)

	res, err := ts.Client.VerifyV2.CreateVerificationCheck(
		ts.ServiceID,
		params,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	output := &OTPServiceVerifyOutput{
		SID:        *res.Sid,
		ServiceSID: *res.ServiceSid,
		AccountSID: *res.AccountSid,
		To:         *res.To,
		Channel:    *res.Channel,
		Status:     *res.Status,
		Valid:      *res.Valid,
	}

	return output, nil
}
