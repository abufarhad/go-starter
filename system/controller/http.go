package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/errors"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/utils/msgutil"
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
func (sys *system) Root(c *gin.Context) {
	c.JSON(http.StatusOK, msgutil.NewRestResp("app backend! let's play!!", nil))
}

// Health will let you know the heart beats ❤️
func (sys *system) Health(c *gin.Context) {
	resp, err := sys.svc.GetHealth()
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", resp), err)
		c.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("Health response", resp))
}
