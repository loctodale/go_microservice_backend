package implement

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_microservice_backend_api/global"
	_const "go_microservice_backend_api/internal/const"
	"go_microservice_backend_api/internal/database"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/utils"
	"go_microservice_backend_api/internal/utils/crypto"
	"go_microservice_backend_api/internal/utils/random"
	"go_microservice_backend_api/internal/utils/sendto"
	"go_microservice_backend_api/pkg/response"
	"log"
	"strconv"
	"strings"
	"time"
)

type sUserLogin struct {
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

func (s *sUserLogin) Login(ctx context.Context) error {
	return nil
}

func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error) {
	//1. hash email
	fmt.Printf("VerifyKey: %s \n", in.VerifyKey)
	fmt.Printf("VerifyType: %s \n", in.VerifyType)

	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("HashKey %s \n", hashKey)

	//2. check user exists in userbase
	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrorCodeUserHasExited, err
	}
	if userFound > 0 {
		return response.ErrorCodeUserHasExited, fmt.Errorf("user %s already exists", in.VerifyKey)
	}
	//3. Create OTP
	userKey := utils.GetUserKey(hashKey)
	// Note: Trong redis nếu một key chưa tồn tại thì return về null chứ không phải có lỗi
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
		break
	case err != nil:
		fmt.Println("Get failed: ", err)
		return response.ErrInvalidOTP, err
		break
	case otpFound != "":
		return response.ErrCodeOTPNotExists, err
	}

	//4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}

	fmt.Printf("OTP is:: %d\n", otpNew)
	// 5. Save OTP in Redis with expired time is 2 minutes
	err = global.Rdb.Set(ctx, userKey, strconv.Itoa(otpNew), time.Duration(_const.TIME_OTP_REGISTER)*time.Minute).Err()
	//
	if err != nil {
		return response.ErrSetRedis, nil
	}
	// 6. Send OTP
	switch in.VerifyType {
	case _const.EMAIL:
		err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, _const.HOST_EMAIL, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrSendEmailOTP, err
		}
		// 7. save OTP to mysql
		result, err := s.r.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
			VerifyKey:     in.VerifyKey,
			VerifyKeyHash: hashKey,
		})
		if err != nil {
			return response.ErrSendEmailOTP, err
		}

		lastIdVerifyUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrSendEmailOTP, err
		}
		log.Println("last verify id user: ", lastIdVerifyUser)
		return response.CodeSuccess, nil
		break
	case _const.MOBILE:
		return response.CodeSuccess, nil
	}

	return response.CodeSuccess, nil
}

func (s *sUserLogin) VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error) {
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	//get otp
	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}
	if in.VerifyCode != otpFound {
		return out, fmt.Errorf("verify code not match")
	}
	infoOTP, err := s.r.GetInfoOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}
	//update status verify
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, err
	}
	out.Token = infoOTP.VerifyKeyHash
	out.Message = "success"
	return out, nil
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context) error {
	return nil
}
