package domain

import "time"

type Post struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Cover      string    `json:"cover"`
	Content    string    `json:"content"`
	UserId     int       `json:"user_id"`
	ClubId     int       `json:"club_id"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
