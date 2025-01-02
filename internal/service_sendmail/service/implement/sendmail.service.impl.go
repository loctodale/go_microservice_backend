package implement

import (
	"go_microservice_backend_api/internal/utils/sendto"
)

type sSendMail struct {
}

func NewSendMailService() *sSendMail {
	return &sSendMail{}
}

func (sm *sSendMail) SendMailVerifyOTP(to []string, content string) error {
	err := sendto.SendTextEmailOtp(to, "thang336655@gmail.com", content)
	if err != nil {
		return err
	}

	return nil
}
