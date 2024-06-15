package service

import (
	"github.com/abufarhad/golang-starter-rest-api/clubs/repository"
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
)

type ClubService struct {
	repo repository.IClubRepo
}

type IClubService interface {
	CreateClub(club dto.ClubReq) (*domain.Club, *errors.RestErr)
	UpdateClub(clubID int, club dto.UpdateClubReq) (*domain.Club, *errors.RestErr)
	GetClubByClubId(clubID int) (*domain.Club, *errors.RestErr)
	DeleteClub(clubID int) *errors.RestErr
	GetAllClubs() ([]domain.Club, *errors.RestErr)
}

func NewClubService(clubRepo repository.IClubRepo) IClubService {
	return &ClubService{
		repo: clubRepo,
	}
}

func (clubSvc *ClubService) CreateClub(club dto.ClubReq) (*domain.Club, *errors.RestErr) {
	newClub := domain.Club{}
	_ = methodutil.StructToStruct(club, &newClub)

	resp, err := clubSvc.repo.CreateClubRepo(newClub)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	return resp, nil
}

func (clubSvc *ClubService) UpdateClub(clubID int, club dto.UpdateClubReq) (*domain.Club, *errors.RestErr) {
	updatedClub := domain.Club{}
	_ = methodutil.StructToStruct(club, &updatedClub)

	updatedClub.Id = clubID

	resp, err := clubSvc.repo.UpdateClubById(updatedClub)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	return resp, nil
}

func (clubSvc *ClubService) GetClubByClubId(clubID int) (*domain.Club, *errors.RestErr) {
	resp, err := clubSvc.repo.GetClubById(clubID)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return resp, nil

}

func (clubSvc *ClubService) DeleteClub(clubID int) *errors.RestErr {
	_, err := clubSvc.repo.GetClubById(clubID)
	if err != nil {
		if errors.As(err, &errors.ErrRecordNotFound) {
			return errors.NewNotFoundError(err)
		}
		return errors.NewInternalServerError(err)
	}

	if err := clubSvc.repo.DeleteClubByID(clubID); err != nil {
		return errors.NewInternalServerError(err)
	}

	return nil
}

func (clubSvc *ClubService) GetAllClubs() ([]domain.Club, *errors.RestErr) {
	resp, err := clubSvc.repo.GetAllClubs()
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return resp, nil
}
