package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/config"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/consts"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
	"github.com/abufarhad/golang-starter-rest-api/user/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IAuth interface {
	Login(req *dto.LoginReq) (*dto.LoginResp, error)
	RefreshToken(refreshToken string) (*dto.LoginResp, error)
}

type auth struct {
	uRepo repository.UserRepo
}

func NewAuthService(uRepo repository.UserRepo) IAuth {
	return &auth{
		uRepo: uRepo,
	}
}

func (as *auth) Login(req *dto.LoginReq) (*dto.LoginResp, error) {
	var user *domain.User
	var err error

	if user, err = as.uRepo.GetUserByEmail(req.Email); err != nil {
		return nil, errors.ErrInvalidEmailOrPassword
	}

	loginPass := []byte(req.Password)
	hashedPass := []byte(*user.Password)

	if err = bcrypt.CompareHashAndPassword(hashedPass, loginPass); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrInvalidEmailOrPassword
	}

	var token *dto.JwtToken

	if token, err = as.CreateToken(uint(user.Id)); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	res := &dto.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return res, nil
}

func (as *auth) RefreshToken(refreshToken string) (*dto.LoginResp, error) {
	oldToken, err := as.parseToken(refreshToken, consts.RefreshTokenType)
	if err != nil {
		return nil, errors.ErrInvalidRefreshToken
	}

	var newToken *dto.JwtToken

	if newToken, err = as.CreateToken(oldToken.UserID); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrCreateJwt
	}

	res := &dto.LoginResp{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}

	return res, nil
}

func (as *auth) parseToken(token, tokenType string) (*dto.JwtToken, error) {
	claims, err := as.parseTokenClaim(token, tokenType)
	if err != nil {
		return nil, err
	}

	tokenDetails := &dto.JwtToken{}

	if err := methodutil.MapToStruct(claims, &tokenDetails); err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}

	if tokenDetails.UserID == 0 || tokenDetails.AccessUuid == "" || tokenDetails.RefreshUuid == "" {
		logger.Info(fmt.Sprintf("%v", claims))
		return nil, errors.ErrInvalidRefreshToken
	}

	return tokenDetails, nil
}

func (as *auth) parseTokenClaim(token, tokenType string) (jwt.MapClaims, error) {
	secret := config.Jwt().AccessTokenSecret

	if tokenType == consts.RefreshTokenType {
		secret = config.Jwt().RefreshTokenSecret
	}

	parsedToken, err := methodutil.ParseJwtToken(token, secret)
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
		return nil, errors.ErrInvalidAccessToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	return claims, nil
}

func (as *auth) CreateToken(userID uint) (*dto.JwtToken, error) {
	var err error
	jwtConf := config.Jwt()
	token := &dto.JwtToken{}

	token.UserID = userID
	token.AccessExpiry = time.Now().Add(time.Minute * jwtConf.AccessTokenExpiry).Unix()
	token.AccessUuid = uuid.New().String()

	token.RefreshExpiry = time.Now().Add(time.Minute * jwtConf.RefreshTokenExpiry).Unix()
	token.RefreshUuid = uuid.New().String()

	user, getErr := as.uRepo.GetUserById(int(userID))
	if getErr != nil {
		return nil, errors.NewError(getErr.Message)
	}

	atClaims := jwt.MapClaims{}
	atClaims["uid"] = user.Id
	atClaims["aid"] = token.AccessUuid
	atClaims["rid"] = token.RefreshUuid
	atClaims["exp"] = token.AccessExpiry

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(jwtConf.AccessTokenSecret))
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrAccessTokenSign
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = user.Id
	rtClaims["aid"] = token.AccessUuid
	rtClaims["rid"] = token.RefreshUuid
	rtClaims["exp"] = token.RefreshExpiry

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(jwtConf.RefreshTokenSecret))
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.ErrRefreshTokenSign
	}

	return token, nil
}
