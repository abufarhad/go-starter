package service

import (
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
	"github.com/abufarhad/golang-starter-rest-api/user/repository"
)

type User struct {
	repo repository.UserRepo
}

type UserService interface {
	CreateUser(user dto.UserReq) (*domain.User, *errors.RestErr)
	UpdateUser(userID int, user dto.UpdateUserReq) (*domain.User, *errors.RestErr)
}

func NewUserService(cusRepo repository.UserRepo) UserService {
	return &User{
		repo: cusRepo,
	}
}

func (usrSvc *User) CreateUser(user dto.UserReq) (*domain.User, *errors.RestErr) {
	usr := domain.User{}
	_ = methodutil.StructToStruct(user, &usr)

	hashPass, _ := methodutil.GenerateHash(usr.Password)
	usr.Password = hashPass
	resp, saveErr := usrSvc.repo.CreateUserRepo(usr)
	if saveErr != nil {
		return nil, saveErr
	}

	return resp, nil
}

func (usrSvc *User) UpdateUser(userID int, usr dto.UpdateUserReq) (*domain.User, *errors.RestErr) {
	user := domain.User{}
	_ = methodutil.StructToStruct(usr, &user)

	user.Id = userID

	updateUser, updateErr := usrSvc.repo.UpdateUserById(user)
	if updateErr != nil {
		return nil, updateErr
	}
	return updateUser, nil
}
