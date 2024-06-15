package dto

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ClubReq struct {
	Name    string `json:"name"`
	OwnerId int    `json:"owner_id"`
	Tags    string `json:"tag,omitempty"`
}

type UpdateClubReq struct {
	Name string `json:"name"`
	Tags string `json:"tag,omitempty"`
}

type ClubResp struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerId   int       `json:"owner_id"`
	Tags      string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

func (vr ClubReq) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Name, v.Required, is.Alphanumeric),
		v.Field(&vr.OwnerId, v.Required, is.Int),
	)
}
