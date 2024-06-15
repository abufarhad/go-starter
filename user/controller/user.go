package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/config"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/msgutil"
	"github.com/abufarhad/golang-starter-rest-api/user/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type User struct {
	svc service.UserService
}

func NewUserController(g interface{}, cusSvc service.UserService) {
	usr := &User{
		svc: cusSvc,
	}
	grp := g.(*gin.RouterGroup)
	grp.GET("/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	grp.POST("/v1/users/signup", usr.CreateUser)
	grp.POST("/v1/users/update", usr.UpdateUser)
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param user body dto.UserReq true "User"
// @Success 200 {object} dto.UserResp{}
// @Failure 400 {object} errors.RestErr
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/users/signup [post]
func (ctr *User) CreateUser(c *gin.Context) {
	reqBody := dto.UserReq{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// validate request
	if payloadErr := reqBody.Validate(); payloadErr != nil {
		logger.ErrorAsJson("failed to validate request body", payloadErr)
		restErr := errors.NewBadRequestError(errors.NewError(errors.ErrRecordNotValid))
		c.JSON(restErr.Status, restErr)
		return
	}

	resp, err := ctr.svc.CreateUser(reqBody)
	if err != nil {
		logger.ErrorAsJson("error", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	userResp := dto.UserResp{}
	_ = methodutil.StructToStruct(*resp, &userResp)
	c.JSON(http.StatusOK, msgutil.NewRestResp("User response", userResp))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param user body dto.UpdateUserReq true "User"
// @Success 200 {object} domain.User{}
// @Failure 400 {object} errors.RestErr
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/users/update [patch]
func (ctr *User) UpdateUser(c *gin.Context) {
	loggedInUser, err := GetUserFromContext(c)
	if err != nil {
		restErr := errors.NewUnauthorizedError(errors.NewError("no logged-in auth found"))
		c.JSON(restErr.Status, restErr)
		return
	}

	var user dto.UpdateUserReq
	if err := c.Bind(&user); err != nil {
		restErr := errors.NewBadRequestError(errors.NewError("invalid json body"))
		c.JSON(restErr.Status, restErr)
		return
	}

	updateUser, updateErr := ctr.svc.UpdateUser(loggedInUser.ID, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("User response", updateUser))
}

func GetUserFromContext(c *gin.Context) (*dto.LoggedInUser, error) {
	userInterface, exists := c.Get(config.Jwt().ContextKey)
	if !exists {
		return nil, errors.ErrNoContextUser
	}

	user, ok := userInterface.(*dto.LoggedInUser)
	if !ok {
		return nil, errors.ErrInvalidUserType
	}

	logger.Info(fmt.Sprintf("%+v", user))

	return user, nil
}
