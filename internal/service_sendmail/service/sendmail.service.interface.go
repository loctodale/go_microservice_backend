package service

type (
	ISendMailService interface {
		SendMailVerifyOTP(to []string, content string) error
	}
)

var localSendMailService ISendMailService

func SendMailService() ISendMailService {
	if localSendMailService == nil {
		panic("please init localSendMailService first")
	}
	return localSendMailService
}

func InitSendMailService(ISendMailService ISendMailService) {
	localSendMailService = ISendMailService
}
