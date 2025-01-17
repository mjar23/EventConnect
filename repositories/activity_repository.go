package repositories

import (
	"database/sql"
	"event-connect/models"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// *************************** ActivityRepository ***************************

// ActivityRepository represents the repository for activity-related database operations
type ActivityRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewActivityRepository creates a new instance of ActivityRepository
func NewActivityRepository(db *sql.DB, logger *logrus.Logger) *ActivityRepository {
	return &ActivityRepository{db: db, logger: logger}
}

// *************************** Repository Methods ***************************

// CreateActivity creates a new activity record in the database
func (r *ActivityRepository) CreateActivity(activity *models.Activity) error {
	// Check if the combination of UserID and EventID already exists in the activities table
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM activities WHERE user_id = $1 AND event_id = $2", activity.UserID, activity.EventID).Scan(&count)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID":       activity.UserID,
			"eventID":      activity.EventID,
			"method":       "CreateActivity",
			"activityType": activity.ActivityType,
		}).Error("Error checking existing activity", err)
		return err
	}

	if count > 0 {
		r.logger.WithFields(logrus.Fields{
			"userID":       activity.UserID,
			"eventID":      activity.EventID,
			"method":       "CreateActivity",
			"activityType": activity.ActivityType,
		}).Warn("Activity already exists for user and event")
		return fmt.Errorf("activity already exists")
	}

	// Insert the new activity record
	query := "INSERT INTO activities (user_id, event_id, activity_type, timestamp) VALUES ($1, $2, $3, $4)"
	_, err = r.db.Exec(query, activity.UserID, activity.EventID, activity.ActivityType, activity.Timestamp)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID":       activity.UserID,
			"eventID":      activity.EventID,
			"method":       "CreateActivity",
			"activityType": activity.ActivityType,
		}).Error("Error executing SQL query", err)
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"userID":       activity.UserID,
		"eventID":      activity.EventID,
		"method":       "CreateActivity",
		"activityType": activity.ActivityType,
	}).Info("Activity created successfully for user and event")
	return nil
}

// GetUserActivities retrieves the activities of a user from the database
func (r *ActivityRepository) GetUserActivities(userID uint) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
        SELECT id, user_id, event_id, activity_type, timestamp
        FROM activities
        WHERE user_id = $1
    `, userID)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID": userID,
			"method": "GetUserActivities",
		}).Error("Error querying user activities", err)
		return nil, err
	}
	defer rows.Close()

	var activities []map[string]interface{}
	for rows.Next() {
		var activity map[string]interface{}
		var id, userID, eventID uint
		var activityType string
		var timestamp time.Time
		err := rows.Scan(&id, &userID, &eventID, &activityType, &timestamp)
		if err != nil {
			r.logger.WithFields(logrus.Fields{
				"userID": userID,
				"method": "GetUserActivities",
			}).Error("Error scanning user activity", err)
			return nil, err
		}
		activity = map[string]interface{}{
			"id":            id,
			"user_id":       userID,
			"event_id":      eventID,
			"activity_type": activityType,
			"timestamp":     timestamp,
		}
		activities = append(activities, activity)
	}
	r.logger.WithFields(logrus.Fields{
		"userID":     userID,
		"method":     "GetUserActivities",
		"activities": activities,
	}).Info("User activities retrieved successfully")
	return activities, nil
}

// GetUserLocationsForEvent retrieves the user locations for a specific event from the database
func (r *ActivityRepository) GetUserLocationsForEvent(eventID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
        SELECT u.id, u.latitude, u.longitude
        FROM activities a
        JOIN users u ON a.user_id = u.id
        WHERE a.event_id = $1 AND a.activity_type = 'event_registered'
    `, eventID)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"eventID": eventID,
			"method":  "GetUserLocationsForEvent",
		}).Error("Error fetching user locations for event", err)
		return nil, err
	}
	defer rows.Close()

	var userLocations []map[string]interface{}
	for rows.Next() {
		var id uint
		var latitude, longitude float64
		err := rows.Scan(&id, &latitude, &longitude)
		if err != nil {
			r.logger.WithFields(logrus.Fields{
				"eventID": eventID,
				"method":  "GetUserLocationsForEvent",
			}).Error("Error scanning user location", err)
			return nil, err
		}
		userLocations = append(userLocations, map[string]interface{}{
			"id":        id,
			"latitude":  latitude,
			"longitude": longitude,
		})
	}

	r.logger.WithFields(logrus.Fields{
		"eventID":      eventID,
		"method":       "GetUserLocationsForEvent",
		"userLocations": userLocations,
	}).Info("User locations for event retrieved successfully")
	return userLocations, nil
}