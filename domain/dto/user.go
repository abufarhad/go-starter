package dto

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type UserReq struct {
	Name     string  `json:"name"`
	Email    string  `gorm:"unique" json:"email"`
	Password *string `json:"password,omitempty"`
}

type UpdateUserReq struct {
	Name string `json:"name"`
}

type UserResp struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (vr UserReq) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Email, v.Required, is.Email),
		v.Field(&vr.Password, v.Required, v.Length(4, 10)),
	)
}
