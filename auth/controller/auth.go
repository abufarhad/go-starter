package controller

import (
	"github.com/abufarhad/golang-starter-rest-api/auth/service"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/msgutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type auth struct {
	authSvc service.IAuth
}

func NewAuthController(g interface{}, authSvc service.IAuth) {
	ath := &auth{
		authSvc: authSvc,
	}
	grp := g.(*gin.RouterGroup)

	grp.POST("/v1/login", ath.Login)
	grp.POST("/v1/token/refresh", ath.RefreshToken)
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param user body dto.LoginReq true "User"
// @Success 200 {object} dto.LoginResp{}
// @Failure 400 {object} errors.RestErr
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/login [post]
func (ctr *auth) Login(c *gin.Context) {
	var cred *dto.LoginReq
	var resp *dto.LoginResp
	var err error

	if err = c.Bind(&cred); err != nil {
		bodyErr := errors.NewBadRequestError(errors.NewError("failed to parse request body"))
		logger.ErrorAsJson("failed to parse request body", err)
		c.JSON(bodyErr.Status, bodyErr)
		return
	}
	cred.Email = strings.TrimSpace(cred.Email)
	cred.Password = strings.TrimSpace(cred.Password)

	if payloadErr := cred.Validate(); payloadErr != nil {
		logger.ErrorAsJson("failed to validate request body", payloadErr)
		restErr := errors.NewBadRequestError(errors.NewError(errors.ErrRecordNotValid))
		c.JSON(restErr.Status, restErr)
		return
	}

	if resp, err = ctr.authSvc.Login(cred); err != nil {
		switch err {
		case errors.ErrInvalidEmailOrPassword:
			unAuthErr := errors.NewUnauthorizedError(errors.ErrInvalidEmailOrPassword)
			c.JSON(unAuthErr.Status, unAuthErr)
			return
		case errors.ErrCreateJwt:
			serverErr := errors.NewInternalServerError(errors.ErrCreateJwt)
			c.JSON(serverErr.Status, serverErr)
			return
		default:
			serverErr := errors.NewInternalServerError(errors.NewError(errors.ErrSomethingWentWrong))
			c.JSON(serverErr.Status, serverErr)
			return
		}
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("auth response", resp))
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh token
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param user body dto.TokenRefreshReq true "User"
// @Success 200 {object} dto.LoginResp{}
// @Failure 400 {object} errors.RestErr
// @Router /api/v1/token/refresh [post]
func (ctr *auth) RefreshToken(c *gin.Context) {
	var token *dto.TokenRefreshReq
	var res *dto.LoginResp
	var err error

	if err = c.Bind(&token); err != nil {
		bodyErr := errors.NewBadRequestError(errors.NewError("failed to parse request body"))
		c.JSON(bodyErr.Status, bodyErr)
		return
	}

	//logger.InfoAsJson("token", token)

	if res, err = ctr.authSvc.RefreshToken(token.RefreshToken); err != nil {
		switch err {
		case errors.ErrParseJwt, errors.ErrInvalidRefreshToken, errors.ErrInvalidRefreshUuid:
			unAuthErr := errors.NewUnauthorizedError(errors.ErrInvalidRefreshToken)
			c.JSON(unAuthErr.Status, unAuthErr)
			return
		case errors.ErrCreateJwt:
			serverErr := errors.NewInternalServerError(errors.ErrCreateJwt)
			c.JSON(serverErr.Status, serverErr)
			return
		default:
			serverErr := errors.NewInternalServerError(errors.NewError("something went wrong"))
			c.JSON(serverErr.Status, serverErr)
			return
		}
	}

	c.JSON(http.StatusOK, msgutil.NewRestResp("refresh_token response", res))
}
