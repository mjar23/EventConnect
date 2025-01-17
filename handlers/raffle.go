package handlers

import (
	"encoding/json"
	"event-connect/models"
	"event-connect/repositories"
	"log"
	"net/http"
	"strconv"
)

func EnterRaffle(raffleRepo *repositories.RaffleRepository, getUserIDFromToken func(*http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			EventID   string  `json:"eventId"`
			Age       int     `json:"age"`
			Gender    string  `json:"gender"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		eventID, err := strconv.ParseUint(requestBody.EventID, 10, 64)
		if err != nil {
			log.Printf("Error parsing eventId: %v", err)
			http.Error(w, "Invalid eventId", http.StatusBadRequest)
			return
		}

		raffleEntry := &models.RaffleEntry{
			EventID:   uint(eventID),
			Age:       requestBody.Age,
			Gender:    requestBody.Gender,
			Latitude:  requestBody.Latitude,
			Longitude: requestBody.Longitude,
		}

		if err := raffleRepo.EnterRaffle(raffleEntry, getUserIDFromToken, r); err != nil {
			log.Printf("Error entering raffle: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
