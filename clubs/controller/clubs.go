package controller

import (
	"github.com/abufarhad/golang-starter-rest-api/user/controller"
	"net/http"
	"strconv"

	"github.com/abufarhad/golang-starter-rest-api/clubs/service"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/msgutil"
	"github.com/gin-gonic/gin"
)

type ClubController struct {
	svc service.IClubService
}

func NewClubController(g interface{}, clubSvc service.IClubService) {
	ctrl := &ClubController{
		svc: clubSvc,
	}
	grp := g.(*gin.RouterGroup)
	grp.POST("/v1/club/create", ctrl.CreateClub)
	grp.PATCH("/v1/club/update", ctrl.UpdateClub)
	grp.GET("/v1/club/:club_id", ctrl.GetClubByClubId)
	grp.DELETE("/v1/club/delete/:club_id", ctrl.DeleteClub)
	grp.GET("/v1/clubs", ctrl.ListOfClubs)
}

// CreateClub godoc
// @Summary Create club
// @Description Create club
// @Tags club
// @Accept application/json
// @Produce application/json
// @Param club body dto.ClubReq true "Club"
// @Success 200 {object} dto.ClubResp{}
// @Failure 400 {object} errors.RestErr
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/club/create [post]

func (ctrl *ClubController) CreateClub(c *gin.Context) {
	reqBody := dto.ClubReq{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	loggedInUser, err := controller.GetUserFromContext(c)
	if err != nil {
		restErr := errors.NewUnauthorizedError(errors.NewError("no logged-in auth found"))
		c.JSON(restErr.Status, restErr)
		return
	}
	reqBody.OwnerId = loggedInUser.ID

	// validate request
	if payloadErr := reqBody.Validate(); payloadErr != nil {
		logger.ErrorAsJson("failed to validate request body", payloadErr)
		restErr := errors.NewBadRequestError(errors.NewError(errors.ErrRecordNotValid))
		c.JSON(restErr.Status, restErr)
		return
	}

	resp, createCluberr := ctrl.svc.CreateClub(reqBody)
	if createCluberr != nil {
		c.JSON(createCluberr.Status, createCluberr)
		return
	}
	clubResp := dto.ClubResp{}
	_ = methodutil.StructToStruct(*resp, &clubResp)
	c.JSON(http.StatusOK, msgutil.NewRestResp("Club response", clubResp))
}

func (ctrl *ClubController) UpdateClub(c *gin.Context) {

	var club dto.UpdateClubReq
	if err := c.Bind(&club); err != nil {
		restErr := errors.NewBadRequestError(errors.NewError("invalid json body"))
		c.JSON(restErr.Status, restErr)
		return
	}

	clubId := c.Param("club_id")
	id, err := strconv.Atoi(clubId)
	if err != nil {
		idErr := errors.NewBadRequestError(errors.NewError("Invalid Club ID"))
		c.JSON(idErr.Status, idErr)
	}

	updatedClub, updateErr := ctrl.svc.UpdateClub(id, club)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("Club response", updatedClub))
}

func (ctrl *ClubController) GetClubByClubId(c *gin.Context) {
	clubID := c.Param("club_id")

	id, err := strconv.Atoi(clubID)
	if err != nil {
		restErr := errors.NewBadRequestError(errors.NewError("Invalid club ID"))
		c.JSON(restErr.Status, restErr)
		return
	}
	club, getClubErr := ctrl.svc.GetClubByClubId(id)
	if getClubErr != nil {
		c.JSON(getClubErr.Status, getClubErr)
		return
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("Club response", club))
}

func (ctrl *ClubController) DeleteClub(c *gin.Context) {
	clubID := c.Param("club_id")
	id, err := strconv.Atoi(clubID)
	if err != nil {
		restErr := errors.NewBadRequestError(errors.NewError("Invalid club ID"))
		c.JSON(restErr.Status, restErr)
		return
	}
	if err := ctrl.svc.DeleteClub(id); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusNoContent, msgutil.NewRestResp("Club deleted successfully", nil))
}

// ListOfClubs godoc
// @Summary List of clubs
// @Description List of clubs
// @Tags club
// @Accept application/json
// @Produce application/json
// @Success 200 {object} dto.ClubResp{}
// @Failure 400 {object} errors.RestErr
// @Failure 500 {object} errors.RestErr
// @Router /api/v1/clubs [get]

func (ctrl *ClubController) ListOfClubs(c *gin.Context) {
	clubs, err := ctrl.svc.GetAllClubs()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, msgutil.NewRestResp("Club Response", clubs))
}
