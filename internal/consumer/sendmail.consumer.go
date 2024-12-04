package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/service"
	"log"
	"strconv"
)

type SendMailConsumer struct {
}
type SendOTPMessage struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

func (sm *SendMailConsumer) NewSendMailConsumer() {
	defer global.KafkaConsumer.Close()
	for {
		m, err := global.KafkaConsumer.ReadMessage(context.Background())

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
		if err = global.KafkaConsumer.CommitMessages(context.Background(), m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
