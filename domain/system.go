package domain

import (
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
)

type ISystemService interface {
	GetHealth() (*dto.HealthResp, *errors.RestErr)
}

type ISystemRepo interface {
	DBCheck() (bool, *errors.RestErr)
}
