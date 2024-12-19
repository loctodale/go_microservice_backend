package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	_const "go_microservice_backend_api/internal/const"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_shop/database"
	"go_microservice_backend_api/internal/service_shop/local"
	"go_microservice_backend_api/internal/utils"
	"go_microservice_backend_api/internal/utils/auth"
	"go_microservice_backend_api/internal/utils/crypto"
	"go_microservice_backend_api/internal/utils/random"
	"go_microservice_backend_api/pkg/response"
	"strconv"
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
		err = global.Rdb.Set(ctx, userKey, otpNew, time.Duration(_const.TIME_OTP_REGISTER)*time.Minute).Err()
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

func (s *sShopRegisterImpl) VerifyOTP(ctx context.Context, in model.VerifyInput) (out model.ShopLoginOutput, err error) {
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}
	if otpFound != in.VerifyCode {
		return out, fmt.Errorf("Token is not valid")
	}

	result, err := s.r.AddIntoShopBase(ctx, database.AddIntoShopBaseParams{
		ShopAccount:  in.VerifyKey,
		ShopPassword: "1234",
	})
	if err != nil {
		return out, err
	}
	//privateKey, _ := auth.GenerateJWTSecret()

	shopId, err := result.LastInsertId()
	if err != nil {
		return out, err
	}

	out.AccessToken, err = auth.NewJWTService().GenerateToken(strconv.FormatInt(shopId, 10), "shop")
	out.RefreshToken = auth.NewJWTService().GenerateRefreshToken(strconv.FormatInt(shopId, 10))
	if err != nil {
		return out, err
	}
	result, err = s.r.AddKeyToken(ctx, database.AddKeyTokenParams{
		ShopID:       uint64(shopId),
		RefreshToken: out.RefreshToken,
	})
	if err != nil {
		return out, err
	}

	out.ShopId = strconv.FormatInt(shopId, 10)
	out.Message = "Verify success, please login and change information"
	return out, nil
}

func (s *sShopRegisterImpl) ChangePasswordRegister(ctx context.Context) (string, error) {
	return "Change password success", nil
}
