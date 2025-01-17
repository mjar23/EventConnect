package handlers

import (
	"database/sql"
	"encoding/json"
	"event-connect/auth"

	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:admin@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

}
func CreateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var comment struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserIDFromToken(r)
	if err != nil {
		log.Printf("Error getting user ID from token: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("INSERT INTO comments (event_id, user_id, text) VALUES ($1, $2, $3)", params["eventId"], userID, comment.Text)
	if err != nil {
		log.Printf("Error inserting comment: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Comment created successfully for event ID %s\n", params["eventId"])
	w.WriteHeader(http.StatusCreated)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID := params["eventId"]

	rows, err := db.Query("SELECT c.id, c.user_id, u.username, c.text, c.created_at FROM comments c JOIN users u ON c.user_id = u.id WHERE c.event_id = $1", eventID)
	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var username, text, createdAt string
		if err := rows.Scan(&id, &userID, &username, &text, &createdAt); err != nil {
			log.Printf("Error scanning comment row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comment := map[string]interface{}{
			"id":        id,
			"userId":    userID,
			"username":  username,
			"text":      text,
			"createdAt": createdAt,
		}
		comments = append(comments, comment)
	}

	log.Printf("Retrieved %d comments for event ID %s", len(comments), eventID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
