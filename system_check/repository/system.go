package repository

import (
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
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

func (sys *system) DBCheck() (bool, error) {
	dB, _ := sys.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}
