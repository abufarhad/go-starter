package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/utils/msgutil"
	"github.com/monstar-lab-bd/golang-starter-rest-api/user/service"
	"net/http"
)

type User struct {
	svc service.UserService
}

func NewUserController(g interface{}, cusSvc service.UserService) {
	cus := &User{
		svc: cusSvc,
	}
	grp := g.(*gin.RouterGroup)

	grp.POST("/v1/user/create", cus.CreateUser)
}

func (cus *User) CreateUser(c *gin.Context) {
	reqBody := domain.User{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	resp, err := cus.svc.CreateUser(reqBody)
	if err != nil {
		logger.ErrorAsJson("error", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("User response", resp))
}
