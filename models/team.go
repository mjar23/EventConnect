package models

import "time"

type Team struct {
	ID        string    `json:"id"`
	EventID   uint      `json:"eventId"`
	EventName string    `json:"eventName"`
	Members   []Member  `json:"members"`
	CreatedAt time.Time `json:"createdAt"`
}

type Member struct {
	UserID            uint    `json:"userId"`
	Username          string  `json:"username"`
	Age               int     `json:"age"`
	Gender            string  `json:"gender"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Email             string  `json:"email"`
	InstagramUsername string  `json:"instagramUsername"`
	FacebookUsername  string  `json:"facebookUsername"`
	SnapchatUsername  string  `json:"snapchatUsername"`
}
