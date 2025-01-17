package models

import (
	_ "github.com/lib/pq" // PostgreSQL driver
)

type RaffleEntry struct {
	EventID   uint    `json:"event_id"`
	UserID    uint    `json:"user_id"`
	Age       int     `json:"age"`
	Gender    string  `json:"gender"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
