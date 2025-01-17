package models

import (
	"time"
)

type User struct {
	ID                 uint      `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Password           string    `json:"password,omitempty"`
	FirstName          *string   `json:"firstName"`
	LastName           string    `json:"lastName"`
	Bio                string    `json:"bio"`
	Interests          string    `json:"interests"`
	Location           string    `json:"location"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	Age                int       `json:"age"`
	Gender             string    `json:"gender"`
	AgeMin             int       `json:"ageMin"`
	AgeMax             int       `json:"ageMax"`
	InstagramUsername  string    `json:"instagramUsername"`
	FacebookUsername   string    `json:"facebookUsername"`
	SnapchatUsername   string    `json:"snapchatUsername"`
	DistancePreference int       `json:"distancePreference"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
