package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service_sendmail/local"
	"go_microservice_backend_api/internal/service_sendmail/service"
	"log"
	"strconv"
)

type SendMailConsumer struct {
}
type SendOTPMessage struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

func InitSendMailConsumer() *SendMailConsumer {
	return &SendMailConsumer{}
}

func (sm *SendMailConsumer) NewSendMailConsumer() {
	consumer := local.KafkaSendMailConsumer
	defer func() {
		err := consumer.Close()
		if err != nil {
			global.Logger.Error(err.Error(), zap.Error(err))
		}
	}()
	for {
		m, err := consumer.ReadMessage(context.Background())

		if err != nil {
			panic(err)
		}
		var data SendOTPMessage
		err = json.Unmarshal(m.Value, &data)
		if err != nil {
			panic(err)
		}
		err = service.SendMailService().SendMailVerifyOTP([]string{data.Email}, strconv.Itoa(data.Otp))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if err = consumer.CommitMessages(context.Background(), m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
