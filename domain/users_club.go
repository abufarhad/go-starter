package domain

type UserClub struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	ClubId int `json:"club_id"`
}
