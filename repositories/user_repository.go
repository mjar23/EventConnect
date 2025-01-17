package repositories

import (
	"database/sql"
	"event-connect/models"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// *************************** UserRepository ***************************

// UserRepository represents the repository for user-related database operations
type UserRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

// *************************** Repository Methods ***************************

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, email, first_name, last_name, bio, interests, location, latitude, longitude, age, gender, age_min, age_max, distance_preference, instagram_username, facebook_username, snapchat_username, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)",
		user.Username, user.Password, user.Email, user.FirstName, user.LastName, user.Bio, user.Interests, user.Location, user.Latitude, user.Longitude, user.Age, user.Gender, user.AgeMin, user.AgeMax, user.DistancePreference, user.InstagramUsername, user.FacebookUsername, user.SnapchatUsername, time.Now(), time.Now())
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"username": user.Username,
			"email":    user.Email,
			"method":   "CreateUser",
		}).Error("Error creating user", err)
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
		"method":   "CreateUser",
	}).Info("User created successfully")
	return nil
}

// GetUserByUsernameAndPassword retrieves a user by their username and password from the database
func (r *UserRepository) GetUserByUsernameAndPassword(username, password string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at
              FROM users
              WHERE username = $1`
	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WithFields(logrus.Fields{
				"username": username,
				"method":   "GetUserByUsernameAndPassword",
			}).Warn("User not found")
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithFields(logrus.Fields{
			"username": username,
			"method":   "GetUserByUsernameAndPassword",
		}).Error("Error retrieving user", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"username": username,
			"method":   "GetUserByUsernameAndPassword",
		}).Warn("Invalid password")
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

// GetUserProfile retrieves a user's profile by their ID from the database
func (r *UserRepository) GetUserProfile(userID uint) (*models.User, error) {
	query := `SELECT id, username, email, first_name, last_name, bio, interests, location, latitude, longitude, age, gender, instagram_username, facebook_username, snapchat_username, created_at, updated_at
			  FROM users
			  WHERE id = $1`

	row := r.db.QueryRow(query, userID)
	var user models.User
	var firstName, lastName, bio, interests, location, instagramUsername, facebookUsername, snapchatUsername sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Email, &firstName, &lastName, &bio, &interests, &location, &user.Latitude, &user.Longitude, &user.Age, &user.Gender, &instagramUsername, &facebookUsername, &snapchatUsername, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WithFields(logrus.Fields{
				"userID": userID,
				"method": "GetUserProfile",
			}).Warn("User not found")
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithFields(logrus.Fields{
			"userID": userID,
			"method": "GetUserProfile",
		}).Error("Error retrieving user profile", err)
		return nil, err
	}

	// Assign nullable fields to user struct
	if firstName.Valid {
		user.FirstName = &firstName.String
	}
	user.LastName = lastName.String
	user.Bio = bio.String
	user.Interests = interests.String
	user.Location = location.String
	user.InstagramUsername = instagramUsername.String
	user.FacebookUsername = facebookUsername.String
	user.SnapchatUsername = snapchatUsername.String

	return &user, nil
}

func (r *UserRepository) UpdateUserProfile(user *models.User) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, email = $2, first_name = $3, last_name = $4, bio = $5, interests = $6, location = $7, latitude = $8, longitude = $9, age = $10, gender = $11, updated_at = $12 WHERE id = $13",
		user.Username, user.Email, user.FirstName, user.LastName, user.Bio, user.Interests, user.Location, user.Latitude, user.Longitude, user.Age, user.Gender, time.Now(), user.ID)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID": user.ID,
			"method": "UpdateUserProfile",
		}).Error("Error updating user profile", err)
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"userID": user.ID,
		"method": "UpdateUserProfile",
	}).Info("User profile updated successfully")
	return nil
}

// GetRecommendedUsers retrieves recommended users based on the user's preferences from the database
func (r *UserRepository) GetRecommendedUsers(user *models.User) ([]models.User, error) {
	query := `
		SELECT id, username, age, latitude, longitude
		FROM users
		WHERE id != $1
		AND age BETWEEN $2 AND $3
		ORDER BY sqrt(power(radians($4 - longitude) * cos(radians((latitude + $5) / 2)), 2) + power(radians(latitude - $5), 2)) * 6371
		LIMIT 10
	`
	rows, err := r.db.Query(query, user.ID, user.AgeMin, user.AgeMax, user.Longitude, user.Latitude)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID": user.ID,
			"method": "GetRecommendedUsers",
		}).Error("Error fetching recommended users", err)
		return nil, err
	}
	defer rows.Close()

	var recommendedUsers []models.User
	for rows.Next() {
		var recommendedUser models.User
		err := rows.Scan(&recommendedUser.ID, &recommendedUser.Username, &recommendedUser.Age, &recommendedUser.Latitude, &recommendedUser.Longitude)
		if err != nil {
			r.logger.WithFields(logrus.Fields{
				"userID": user.ID,
				"method": "GetRecommendedUsers",
			}).Error("Error scanning recommended user", err)
			return nil, err
		}
		recommendedUsers = append(recommendedUsers, recommendedUser)
	}

	return recommendedUsers, nil
}

// UpdateUserPreferences updates a user's preferences in the database
func (r *UserRepository) UpdateUserPreferences(userID uint, ageMin, ageMax, distancePreference int) error {
	_, err := r.db.Exec("UPDATE users SET age_min = $1, age_max = $2, distance_preference = $3 WHERE id = $4",
		ageMin, ageMax, distancePreference, userID)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"userID": userID,
			"method": "UpdateUserPreferences",
		}).Error("Error updating user preferences", err)
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"userID": userID,
		"method": "UpdateUserPreferences",
	}).Info("User preferences updated successfully")
	return nil
}

// GetUserByID retrieves a user by their ID from the database
func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	query := `SELECT id, username, email, first_name, last_name, bio, interests, location, latitude, longitude, age, gender, instagram_username, facebook_username, snapchat_username, created_at, updated_at
			  FROM users
			  WHERE id = $1`

	row := r.db.QueryRow(query, userID)
	var user models.User
	var firstName, lastName, bio, interests, location, instagramUsername, facebookUsername, snapchatUsername sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Email, &firstName, &lastName, &bio, &interests, &location, &user.Latitude, &user.Longitude, &user.Age, &user.Gender, &instagramUsername, &facebookUsername, &snapchatUsername, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WithFields(logrus.Fields{
				"userID": userID,
				"method": "GetUserByID",
			}).Warn("User not found")
			return nil, fmt.Errorf("user not found")
		}
		r.logger.WithFields(logrus.Fields{
			"userID": userID,
			"method": "GetUserByID",
		}).Error("Error retrieving user", err)
		return nil, err
	}

	// Assign nullable fields to user struct
	if firstName.Valid {
		user.FirstName = &firstName.String
	}
	user.LastName = lastName.String
	user.Bio = bio.String
	user.Interests = interests.String
	user.Location = location.String
	user.InstagramUsername = instagramUsername.String
	user.FacebookUsername = facebookUsername.String
	user.SnapchatUsername = snapchatUsername.String

	return &user, nil
}