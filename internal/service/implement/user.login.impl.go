package implement

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go_microservice_backend_api/global"
	_const "go_microservice_backend_api/internal/const"
	"go_microservice_backend_api/internal/database"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/utils"
	"go_microservice_backend_api/internal/utils/auth"
	"go_microservice_backend_api/internal/utils/crypto"
	"go_microservice_backend_api/internal/utils/random"
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

func (s *sUserLogin) Login(ctx context.Context, in model.LoginInput) (codeResult int, out model.LoginOutput, err error) {
	userBaseFound, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	//1. check account is exists
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	if !crypto.MatchPassword(userBaseFound.UserPassword, in.UserPassword, userBaseFound.UserSalt) {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("Dose not match password")
	}
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserAccount:  in.UserAccount,
		UserLoginIp:  sql.NullString{String: "127.0.0.1", Valid: true},
		UserPassword: in.UserPassword,
	})
	subToken := utils.GenerateCliTokenUUID(int(userBaseFound.UserID))
	infoUser, err := s.r.GetUser(ctx, uint64(userBaseFound.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	infoUserJson, err := json.Marshal(infoUser)
	if err != nil {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("json marshal err: %v", err)
	}
	err = global.Rdb.Set(ctx, subToken, infoUserJson, _const.TIME_OTP_REGISTER*time.Minute).Err()
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}

	//Generate token
	out.Token, err = auth.CreateToken(subToken)
	if err != nil {
		return
	}
	return 200, out, nil
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
	// Note: Trong redis nếu một certs chưa tồn tại thì return về null chứ không phải có lỗi
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
		//err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, _const.HOST_EMAIL, strconv.Itoa(otpNew))
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
		body := make(map[string]interface{})
		body["email"] = in.VerifyKey
		body["otp"] = otpNew

		bodyRequest, _ := json.Marshal(body)
		message := kafka.Message{
			Key:   []byte("otp-auth"),
			Value: []byte(bodyRequest),
			Time:  time.Now(),
		}
		err = global.KafkaProducer.WriteMessages(context.Background(), message)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
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

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	infoOTP, err := s.r.GetInfoOTP(ctx, token)
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	if infoOTP.IsVerified.Int32 == 0 {
		return response.ErrCodeUSerOTPNotExists, fmt.Errorf("User not verify")
	}

	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = infoOTP.VerifyKey
	userBase.UserSalt, err = crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	userBase.UserPassword = crypto.HashPasswordSalt(password, userBase.UserSalt)
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	// add user_id to user_info
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:               uint64(user_id),
		UserAccount:          infoOTP.VerifyKey,
		UserNickname:         sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserAvatar:           sql.NullString{String: "", Valid: true},
		UserState:            1,
		UserMobile:           sql.NullString{String: "", Valid: true},
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:         sql.NullTime{Time: time.Now(), Valid: true},
		UserEmail:            sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserIsAuthentication: 1,
	})
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	user_id, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrCodeUSerOTPNotExists, err
	}
	return int(user_id), nil
}
