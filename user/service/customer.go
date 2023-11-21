package service

import (
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/errors"
	"github.com/monstar-lab-bd/golang-starter-rest-api/user/repository"
)

type User struct {
	repo repository.UserRepo
}

type UserService interface {
	CreateUser(order domain.User) (*domain.User, *errors.RestErr)
}

func NewUserService(cusRepo repository.UserRepo) UserService {
	return &User{
		repo: cusRepo,
	}
}

func (cus *User) CreateUser(order domain.User) (*domain.User, *errors.RestErr) {

	resp, saveErr := cus.repo.CreateUserRepo(order)
	if saveErr != nil {
		return nil, saveErr
	}

	return resp, nil
}
