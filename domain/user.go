package domain

import "time"

type User struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Password  *string    `json:"password,omitempty"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
