package service

import (
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
)

type system struct {
	repo domain.ISystemRepo
}

func NewSystemService(sysRepo domain.ISystemRepo) domain.ISystemService {
	return &system{
		repo: sysRepo,
	}
}

func (sys *system) GetHealth() (*dto.HealthResp, *errors.RestErr) {
	resp := dto.HealthResp{}

	// check db
	dbOnline, err := sys.repo.DBCheck()
	resp.DBOnline = dbOnline

	if err != nil {
		return &resp, err
	}

	return &resp, nil
}
