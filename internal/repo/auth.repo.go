package repo

import (
	"fmt"
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/database"
	"time"
)

type IAuthRepository interface {
	AddOTP(email string, otp int, expiredTime int64) error
}

type authRepository struct {
	sqlc *database.Queries
}

func (ar *authRepository) AddOTP(email string, otp int, expiredTime int64) error {
	key := fmt.Sprintf("usr:%s:otp", email)
	return global.Rdb.SetEx(Ctx, key, otp, time.Duration(expiredTime)).Err()
}

func NewAuthRepository() IAuthRepository {
	return &authRepository{
		sqlc: database.New(global.Mdb),
	}
}
