package repository

import (
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/errors"
	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
}

type UserRepo interface {
	CreateUserRepo(User domain.User) (*domain.User, *errors.RestErr)
	GetUserById(cID int) (*domain.User, *errors.RestErr)
	UpdateUserInfoById(User domain.User) *errors.RestErr
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &User{
		DB: db,
	}
}

func (cus *User) CreateUserRepo(User domain.User) (*domain.User, *errors.RestErr) {
	err := cus.DB.Model(&domain.User{}).Create(&User).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &User, nil
}

func (cus *User) GetUserById(cID int) (*domain.User, *errors.RestErr) {
	resp := domain.User{}
	err := cus.DB.Model(&domain.User{}).Where("id = ?", cID).Find(&resp).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &resp, nil
}

func (cus *User) UpdateUserInfoById(User domain.User) *errors.RestErr {
	if err := cus.DB.Model(&domain.User{}).Where("id=?", User.Id).Updates(User).Error; err != nil {
		return errors.NewInternalServerError(err)
	}
	return nil
}
