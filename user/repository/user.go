package repository

import (
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
}

type UserRepo interface {
	CreateUserRepo(user domain.User) (*domain.User, *errors.RestErr)
	GetUserById(cID int) (*domain.User, *errors.RestErr)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserById(user domain.User) (*domain.User, *errors.RestErr)
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &User{
		DB: db,
	}
}

func (userRep *User) CreateUserRepo(user domain.User) (*domain.User, *errors.RestErr) {
	err := userRep.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &user, nil
}

func (userRep *User) GetUserById(cID int) (*domain.User, *errors.RestErr) {
	resp := domain.User{}
	err := userRep.DB.Model(&domain.User{}).Where("id = ?", cID).Find(&resp).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if resp.Id == 0 {
		return nil, errors.NewNotFoundError(errors.NewError(errors.ErrRecordNotFound))
	}
	return &resp, nil
}

func (userRep *User) GetUserByEmail(email string) (*domain.User, error) {
	resp := domain.User{}
	err := userRep.DB.Model(&domain.User{}).Where("email = ?", email).Find(&resp).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if resp.Id == 0 {
		return nil, errors.NewNotFoundError(errors.NewError(errors.ErrRecordNotFound))
	}
	return &resp, nil
}

func (userRep *User) UpdateUserById(user domain.User) (*domain.User, *errors.RestErr) {
	if err := userRep.DB.Model(&domain.User{}).Where("id=?", user.Id).Updates(user).Error; err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &user, nil
}
