package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/repo"
	"go_microservice_backend_api/internal/utils/crypto"
	"go_microservice_backend_api/internal/utils/random"
	"go_microservice_backend_api/pkg/response"

	"time"
)

type IUserService interface {
	Register(email string, purpose string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IAuthRepository
}

func NewUserService(userRepo repo.IUserRepository, userAuthRepo repo.IAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

func (us *userService) Register(email string, purpose string) int {
	//0. hash email
	hashEmail := crypto.GetHash(email)
	fmt.Printf("hashEmail:%v\n", hashEmail)
	//1. check exists email in db
	foundUser := us.userRepo.GetUserByEmail(email)
	//fmt.Printf("foundUser:%v\n", foundUser.UsrID)
	if foundUser {
		return response.ErrorCodeUserHasExited
	}
	//2. new otp
	otp := random.GenerateSixDigitOtp()
	if purpose == "TEST_USER" {
		otp = 1234566
	}

	fmt.Printf("Otp is:: %d\n", otp)

	//3. save otp in redis with expiration time
	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))

	if err != nil {
		println(err.Error())
		return response.ErrInvalidOTP
	}
	//4. send email otp
	//err = sendto.SendTextEmailOtp([]string{email}, "thang336655@gmail.com", strconv.Itoa(otp))
	//err = sendto.SendTemplateEmailOTP([]string{email}, "Go-Backend-API", "otp-auth.html", map[string]interface{}{
	//	"otp": strconv.Itoa(otp),
	//})
	//if err != nil {
	//	println("ErrSendEmailOTP")
	//	return response.ErrSendEmailOTP
	//}

	body := make(map[string]interface{})
	body["email"] = email
	body["otp"] = otp

	bodyRequest, _ := json.Marshal(body)
	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}

	err = global.KafkaProducer.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Println(err.Error())
		return response.ErrSendEmailOTP
	}
	//5. check otp available

	//6. user spam

	return response.CodeSuccess
}
