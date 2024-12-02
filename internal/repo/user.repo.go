package repo

import (
	"go_microservice_backend_api/global"
	"go_microservice_backend_api/internal/database"
)

type IUserRepository interface {
	GetUserByEmail(email string) bool
}

type userRepository struct {
	sqlc *database.Queries
}

func (ur *userRepository) GetUserByEmail(email string) bool {
	//data, _ := ur.sqlc.GetUserByEmail(Ctx, email)
	//return data
	return true
}

func NewUserRepository() IUserRepository {
	return &userRepository{
		sqlc: database.New(global.Mdb),
	}
}
