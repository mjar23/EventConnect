package repositories

import (
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(eventID string, userID int, text string) error {
	_, err := r.db.Exec("INSERT INTO comments (event_id, user_id, text) VALUES ($1, $2, $3)", eventID, userID, text)
	return err
}

func (r *CommentRepository) GetComments(eventID string) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.user_id, u.username, c.text, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.event_id = $1
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var username, text, createdAt string
		if err := rows.Scan(&id, &userID, &username, &text, &createdAt); err != nil {
			return nil, err
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

	return comments, nil
}
