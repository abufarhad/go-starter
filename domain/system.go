package domain

import "github.com/monstar-lab-bd/golang-starter-rest-api/domain/dto"

type ISystemService interface {
	GetHealth() (*dto.HealthResp, error)
}

type ISystemRepo interface {
	DBCheck() (bool, error)
}
