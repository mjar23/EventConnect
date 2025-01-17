package handlers

import (
	"encoding/json"
	"event-connect/models"
	"event-connect/repositories"
	"log"
	"net/http"
	"time"
)

func CreateActivity(activityRepo *repositories.ActivityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var activity models.Activity
		err := json.NewDecoder(r.Body).Decode(&activity)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		activity.Timestamp = time.Now()
		err = activityRepo.CreateActivity(&activity)
		if err != nil {
			log.Printf("Error creating activity: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetUserActivities(activityRepo *repositories.ActivityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "Missing userId parameter", http.StatusBadRequest)
			return
		}

		activities, err := activityRepo.GetUserActivities(uint(0)) // Modify to pass the actual user ID
		if err != nil {
			log.Printf("Error fetching user activities: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activities)
	}
}

func GetUserLocationsForEvent(activityRepo *repositories.ActivityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := r.URL.Query().Get("eventId")
		if eventID == "" {
			http.Error(w, "Missing eventId parameter", http.StatusBadRequest)
			return
		}

		userLocations, err := activityRepo.GetUserLocationsForEvent(eventID)
		if err != nil {
			log.Printf("Error fetching user locations for event: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userLocations)
	}
}
