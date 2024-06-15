package dto

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenRefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoggedInUser struct {
	ID          int    `json:"user_id"`
	AccessUuid  string `json:"access_uuid"`
	RefreshUuid string `json:"refresh_uuid"`
}

type JwtToken struct {
	UserID        uint   `json:"uid"`
	AccessToken   string `json:"act"`
	RefreshToken  string `json:"rft"`
	AccessUuid    string `json:"aid"`
	RefreshUuid   string `json:"rid"`
	AccessExpiry  int64  `json:"axp"`
	RefreshExpiry int64  `json:"rxp"`
}

func (vr LoginReq) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Email, v.Required, is.Email),
		v.Field(&vr.Password, v.Required),
	)
}
