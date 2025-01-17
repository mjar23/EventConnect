package models

import (
	"time"
)

type Activity struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	EventID      uint      `json:"event_id"`
	ActivityType string    `json:"activity_type"`
	Timestamp    time.Time `json:"timestamp"`
}
