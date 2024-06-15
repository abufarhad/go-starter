package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/msgutil"
	"net/http"
)

type system struct {
	svc domain.ISystemService
}

// NewSystemController will initialize the controllers
func NewSystemController(g interface{}, sysSvc domain.ISystemService) {
	sc := &system{
		svc: sysSvc,
	}
	grp := g.(*gin.RouterGroup)

	grp.GET("/v1", sc.Root)
	grp.GET("/v1/h34l7h", sc.Health)
}

// Root will let you see what you can slash 🐲
// SystemCheck godoc
// @Summary Root
// @Description Root
// @Tags system
// @Accept application/json
// @Produce application/json
// @Success 200 {object} msgutil.RestResp{}
// @Failure 500 {object} errors.RestErr
// @Router /api/v1 [get]
func (sys *system) Root(c *gin.Context) {
	c.JSON(http.StatusOK, msgutil.NewRestResp("app backend! let's play!!", nil))
}

// Health will let you know the heart beats ❤️
// Health godoc
// @Summary Health
// @Description Health
// @Tags system
// @Accept application/json
// @Produce application/json
// @Success 200 {object} msgutil.RestResp{}
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/h34l7h [get]
func (sys *system) Health(c *gin.Context) {
	resp, err := sys.svc.GetHealth()
	if err != nil {
		logger.ErrorAsJson("error", err)
		c.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("Health response", resp))
}
