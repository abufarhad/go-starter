package repository

import (
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"gorm.io/gorm"
)

type ClubRepo struct {
	*gorm.DB
}

type IClubRepo interface {
	CreateClubRepo(club domain.Club) (*domain.Club, *errors.RestErr)
	GetClubById(cID int) (*domain.Club, *errors.RestErr)
	UpdateClubById(club domain.Club) (*domain.Club, *errors.RestErr)
	DeleteClubByID(clubID int) *errors.RestErr
	GetAllClubs() ([]domain.Club, *errors.RestErr)
}

func NewClubRepository(db *gorm.DB) IClubRepo {
	return &ClubRepo{
		DB: db,
	}
}

func (clubRepo *ClubRepo) CreateClubRepo(club domain.Club) (*domain.Club, *errors.RestErr) {
	err := clubRepo.DB.Model(&domain.Club{}).Create(&club).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &club, nil
}

func (clubRepo *ClubRepo) GetClubById(clubID int) (*domain.Club, *errors.RestErr) {
	resp := domain.Club{}
	err := clubRepo.DB.Model(&domain.Club{}).Where("id=?", clubID).Find(&resp).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	if resp.Id == 0 {
		return nil, errors.NewNotFoundError(errors.NewError(errors.ErrRecordNotFound))
	}
	return &resp, nil
}

func (clubRep *ClubRepo) UpdateClubById(club domain.Club) (*domain.Club, *errors.RestErr) {
	if err := clubRep.DB.Model(&domain.Club{}).Where("id=?", club.Id).Updates(club).Error; err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return &club, nil
}

func (clubRepo *ClubRepo) DeleteClubByID(clubID int) *errors.RestErr {
	if err := clubRepo.DB.Where("id = ?", clubID).Delete(&domain.Club{}).Error; err != nil {
		return errors.NewInternalServerError(err)
	}
	return nil
}

func (clubRepo *ClubRepo) GetAllClubs() ([]domain.Club, *errors.RestErr) {
	var clubs []domain.Club
	err := clubRepo.DB.Model(&domain.Club{}).Find(&clubs).Error
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	return clubs, nil
}
