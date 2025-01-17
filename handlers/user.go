package handlers

import (
	"encoding/json"
	"event-connect/auth"
	"event-connect/models"
	"event-connect/repositories"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Request Body: %+v", user)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		err = userRepo.CreateUser(&user)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Created User: %+v", user)

		token, err := auth.GenerateJWT(user.ID)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.Password = ""

		response := struct {
			User  models.User `json:"user"`
			Token string      `json:"token"`
		}{
			User:  user,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetUserInfo(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromToken(r)
		if err != nil {
			log.Printf("Error getting user ID from token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := userRepo.GetUserByID(uint(userID))
		if err != nil {
			log.Printf("Error fetching user info: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userInfo := struct {
			Age       int     `json:"age"`
			Gender    string  `json:"gender"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Age:       user.Age,
			Gender:    user.Gender,
			Latitude:  user.Latitude,
			Longitude: user.Longitude,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	}
}

func GetUserProfile(userRepo *repositories.UserRepository, activityRepo *repositories.ActivityRepository, teamRepo *repositories.TeamRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := userRepo.GetUserProfile(uint(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		activities, err := activityRepo.GetUserActivities(uint(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teams, err := teamRepo.FetchUserTeams(uint(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"user":       user,
			"activities": activities,
			"teams":      teams,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateUserProfile(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user models.User
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = uint(userID)

		err = userRepo.UpdateUserProfile(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User profile updated successfully"})
	}
}

func GetOtherUserProfile(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		uid, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := userRepo.GetUserProfile(uint(uid))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
