package repository

import (
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"gorm.io/gorm"
)

type system struct {
	*gorm.DB
}

func NewSystemRepository(db *gorm.DB) domain.ISystemRepo {
	return &system{
		DB: db,
	}
}

func (sys *system) DBCheck() (bool, *errors.RestErr) {
	dB, _ := sys.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, errors.NewInternalServerError(err)
	}

	return true, nil
}
