package repositories

import (
	"database/sql"
	"event-connect/models"
	"event-connect/skiddle"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RaffleRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewRaffleRepository(db *sql.DB, logger *logrus.Logger) *RaffleRepository {
	return &RaffleRepository{db: db, logger: logger}
}

func (r *RaffleRepository) EnterRaffle(entry *models.RaffleEntry, getUserIDFromToken func(*http.Request) (int, error), req *http.Request) error {
	url := skiddle.EventDetailsURL(entry.EventID)
	resp, err := http.Get(url)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"method":  "EnterRaffle",
		}).Error("Error checking event existence", err)
		return fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"method":  "EnterRaffle",
		}).Error("Event not found in Skiddle API")
		return fmt.Errorf("event not found")
	}

	userID, err := getUserIDFromToken(req)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"method":  "EnterRaffle",
		}).Error("Error getting user ID from token", err)
		return fmt.Errorf("unauthorized")
	}

	entry.UserID = uint(userID)

	r.logger.WithFields(logrus.Fields{
		"eventID": entry.EventID,
		"userID":  entry.UserID,
		"method":  "EnterRaffle",
	}).Info("Received raffle entry request", entry)

	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM raffle_entries WHERE event_id = $1 AND user_id = $2", entry.EventID, entry.UserID).Scan(&count)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"userID":  entry.UserID,
			"method":  "EnterRaffle",
		}).Error("Error checking duplicate raffle entry", err)
		return fmt.Errorf("internal server error")
	}

	if count > 0 {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"userID":  entry.UserID,
			"method":  "EnterRaffle",
		}).Warn("User has already entered the raffle for event")
		return fmt.Errorf("duplicate raffle entry")
	}

	_, err = r.db.Exec("INSERT INTO raffle_entries (event_id, user_id, age, gender, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6)",
		entry.EventID, entry.UserID, entry.Age, entry.Gender, entry.Latitude, entry.Longitude)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"eventID": entry.EventID,
			"userID":  entry.UserID,
			"method":  "EnterRaffle",
		}).Error("Error inserting raffle entry", err)
		return fmt.Errorf("internal server error")
	}

	r.logger.WithFields(logrus.Fields{
		"eventID": entry.EventID,
		"userID":  entry.UserID,
		"method":  "EnterRaffle",
	}).Info("Raffle entry created successfully")
	return nil
}