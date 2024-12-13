package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_shop/database"
	"go_microservice_backend_api/internal/service_shop/local"
	"go_microservice_backend_api/internal/utils"
	"go_microservice_backend_api/internal/utils/crypto"
	"go_microservice_backend_api/internal/utils/random"
	"go_microservice_backend_api/pkg/response"
	"strings"
	"time"
)

type sShopRegisterImpl struct {
	r *database.Queries
}

func NewShopRegister(r *database.Queries) *sShopRegisterImpl {
	return &sShopRegisterImpl{r: r}
}

func (s *sShopRegisterImpl) Register(ctx context.Context, in model.RegisterInput) (codeResult int, err error) {
	fmt.Printf("VerifyKey: %s \n", in.VerifyKey)
	shopBaseFound, err := s.r.CheckShopBaseIsExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrorCodeUserHasExited, err
	}
	if shopBaseFound > 0 {
		return response.ErrorCodeUserHasExited, fmt.Errorf("shop base is exists")
	}

	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	userKey := utils.GetUserKey(hashKey)
	fmt.Println("Hash key: ", hashKey)
	fmt.Println("UserKey: ", userKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()
	fmt.Println("otpFound: ", otpFound)
	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
		break
	case err != nil:
		fmt.Println("Get failed: ", err)
		return response.ErrInvalidOTP, err
		break
	case otpFound != "":
		return response.ErrCodeOTPNotExists, fmt.Errorf("Error found != ")
	}

	otpNew := random.GenerateSixDigitOtp()
	done := make(chan error)
	go func() {
		err = global.Rdb.Set(ctx, userKey, otpNew, 0).Err()
		if err != nil {
			done <- err
		}
		done <- nil
	}()

	if err = <-done; err != nil {
		fmt.Println("Go routine fail: ", err.Error())
		return response.ErrorCodeUserHasExited, err
	}
	messageBody := make(map[string]interface{})
	messageBody["email"] = in.VerifyKey
	messageBody["otp"] = otpNew

	bodyRequest, _ := json.Marshal(messageBody)
	message := kafka.Message{
		Key:   []byte(userKey),
		Value: bodyRequest,
		Time:  time.Now(),
	}
	err = local.ShopProducer.WriteMessages(ctx, message)
	if err != nil {
		fmt.Println("Producer fail: ", err.Error())
		return response.ErrorCodeUserHasExited, err
	}
	return response.CodeSuccess, nil
}
