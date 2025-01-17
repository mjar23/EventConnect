// handlers/login.go
package handlers

import (
	"encoding/json"
	"event-connect/auth" // Import the auth package correctly
	"event-connect/repositories"
	"net/http"
)

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for the login endpoint
type LoginResponse struct {
	Token string `json:"token"`
}

func Login(userRepo *repositories.UserRepository, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user from the database using the UserRepository
	user, err := userRepo.GetUserByUsernameAndPassword(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
