package repositories

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "math/rand"
    "time"

    "event-connect/models"

    "github.com/sirupsen/logrus"
)

// *************************** TeamRepository ***************************

// TeamRepository represents the repository for team-related database operations
type TeamRepository struct {
    db     *sql.DB
    logger *logrus.Logger
    
}

// NewTeamRepository creates a new instance of TeamRepository
func NewTeamRepository(db *sql.DB, logger *logrus.Logger) *TeamRepository {
    return &TeamRepository{db: db, logger: logger}
}

// *************************** Repository Methods ***************************

// FetchRaffleEntries fetches the raffle entries for a specific event from the database
func (r *TeamRepository) FetchRaffleEntries(eventID uint) ([]models.User, error) {
    rows, err := r.db.Query("SELECT user_id, age, gender, latitude, longitude FROM raffle_entries WHERE event_id = $1 ORDER BY gender, age, latitude, longitude", eventID)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "FetchRaffleEntries",
        }).Error("Failed to fetch raffle entries", err)
        return nil, err
    }
    defer rows.Close()

    var entries []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Age, &user.Gender, &user.Latitude, &user.Longitude)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "eventId": eventID,
                "method":  "FetchRaffleEntries",
            }).Error("Failed to scan raffle entry", err)
            return nil, err
        }
        entries = append(entries, user)
    }

    return entries, nil
}

// InsertTeams inserts the teams into the database for a specific event
func (r *TeamRepository) InsertTeams(eventID uint, teams []models.Team) error {
    tx, err := r.db.Begin()
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "InsertTeams",
        }).Error("Failed to begin transaction", err)
        return err
    }

    for _, team := range teams {
        teamID := generateUniqueTeamID()

        for _, member := range team.Members {
            _, err := tx.Exec("INSERT INTO teams (event_id, user_id, age, gender, latitude, longitude, team_id, email, instagram_username, facebook_username, snapchat_username) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
                eventID, member.UserID, member.Age, member.Gender, member.Latitude, member.Longitude, teamID, member.Email, member.InstagramUsername, member.FacebookUsername, member.SnapchatUsername)
            if err != nil {
                tx.Rollback()
                r.logger.WithFields(logrus.Fields{
                    "eventId": eventID,
                    "method":  "InsertTeams",
                }).Error("Failed to insert team member", err)
                return err
            }

            // Fetch user details, including email and social media usernames
            user, err := r.GetUserByID(member.UserID)
            if err != nil {
                tx.Rollback()
                r.logger.WithFields(logrus.Fields{
                    "eventId": eventID,
                    "method":  "InsertTeams",
                    "userId":  member.UserID,
                }).Error("Failed to get user details", err)
                return err
            }

            member.Email = user.Email
            member.InstagramUsername = user.InstagramUsername
            member.FacebookUsername = user.FacebookUsername
            member.SnapchatUsername = user.SnapchatUsername
        }
    }

    err = tx.Commit()
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "InsertTeams",
        }).Error("Failed to commit transaction", err)
        return err
    }

    return nil
}

// GetUserByID retrieves a user by their ID from the database
func (r *TeamRepository) GetUserByID(userID uint) (*models.User, error) {
    row := r.db.QueryRow("SELECT email, instagram_username, facebook_username, snapchat_username FROM users WHERE id = $1", userID)
    var user models.User
    err := row.Scan(&user.Email, &user.InstagramUsername, &user.FacebookUsername, &user.SnapchatUsername)
    if err != nil {
        if err == sql.ErrNoRows {
            r.logger.WithFields(logrus.Fields{
                "userId": userID,
                "method": "GetUserByID",
            }).Info("User not found")
            return nil, nil
        }
        r.logger.WithFields(logrus.Fields{
            "userId": userID,
            "method": "GetUserByID",
        }).Error("Failed to get user details", err)
        return nil, err
    }
    user.ID = userID
    return &user, nil
}

// FetchUserTeams retrieves the teams of a user from the database
func (r *TeamRepository) FetchUserTeams(userID uint) ([]models.Team, error) {
    rows, err := r.db.Query(`
        SELECT t.event_id, t.team_id, t.created_at,
            json_agg(json_build_object('userId', u.id, 'username', u.username, 'age', u.age, 'gender', u.gender)) AS members
        FROM teams t
        JOIN users u ON t.user_id = u.id
        WHERE t.team_id IN (
            SELECT DISTINCT team_id
            FROM teams
            WHERE team_id IN (
                SELECT team_id
                FROM teams
                WHERE user_id = $1
            )
        )
        GROUP BY t.event_id, t.team_id, t.created_at
    `, userID)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "userId": userID,
            "method": "FetchUserTeams",
        }).Error("Failed to fetch user teams", err)
        return nil, err
    }
    defer rows.Close()

    var teams []models.Team
    for rows.Next() {
        var team models.Team
        var membersJSON string
        err := rows.Scan(&team.EventID, &team.ID, &team.CreatedAt, &membersJSON)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "userId": userID,
                "method": "FetchUserTeams",
            }).Error("Failed to scan user team", err)
            return nil, err
        }

        err = json.Unmarshal([]byte(membersJSON), &team.Members)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "userId": userID,
                "method": "FetchUserTeams",
            }).Error("Failed to unmarshal user team members", err)
            return nil, err
        }

        teams = append(teams, team)
    }

    return teams, nil
}

// GetTeamsForEvent retrieves the teams for a specific event from the database
func (r *TeamRepository) GetTeamsForEvent(eventID uint) ([]models.Team, error) {
    rows, err := r.db.Query(`
        SELECT json_agg(json_build_object('id', team_id, 'members', members)) AS teams
        FROM (
            SELECT team_id, json_agg(json_build_object('userId', user_id, 'age', age, 'gender', gender, 'email', email, 'instagram_username', instagram_username, 'facebook_username', facebook_username, 'snapchat_username', snapchat_username)) AS members
            FROM teams
            WHERE event_id = $1
            GROUP BY team_id
        ) t
    `, eventID)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "GetTeamsForEvent",
        }).Error("Failed to fetch teams for event", err)
        return nil, err
    }
    defer rows.Close()

    var teamsJSON string
    if rows.Next() {
        err := rows.Scan(&teamsJSON)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "eventId": eventID,
                "method":  "GetTeamsForEvent",
            }).Error("Failed to scan teams for event", err)
            return nil, err
        }
    }

    var teams []models.Team
    err = json.Unmarshal([]byte(teamsJSON), &teams)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "GetTeamsForEvent",
        }).Error("Failed to unmarshal teams for event", err)
        return nil, err
    }

    return teams, nil
}

// FetchEventIDsFromRaffleEntries retrieves the distinct event IDs from the raffle entries in the database
func (r *TeamRepository) FetchEventIDsFromRaffleEntries() ([]uint, error) {
    rows, err := r.db.Query(`
        SELECT DISTINCT event_id
        FROM raffle_entries
    `)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "method": "FetchEventIDsFromRaffleEntries",
        }).Error("Failed to fetch event IDs from raffle entries", err)
        return nil, err
    }
    defer rows.Close()

    var eventIDs []uint
    for rows.Next() {
        var eventID uint
        err := rows.Scan(&eventID)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "method": "FetchEventIDsFromRaffleEntries",
            }).Error("Failed to scan event ID from raffle entries", err)
            return nil, err
        }
        eventIDs = append(eventIDs, eventID)
    }

    return eventIDs, nil
}

// FetchRaffleEntriesByEventID retrieves the raffle entries for a specific event from the database
func (r *TeamRepository) FetchRaffleEntriesByEventID(eventID uint) ([]models.User, error) {
    rows, err := r.db.Query(`
        SELECT user_id, age, gender, latitude, longitude
        FROM raffle_entries
        WHERE event_id = $1
    `, eventID)
    if err != nil {
        r.logger.WithFields(logrus.Fields{
            "eventId": eventID,
            "method":  "FetchRaffleEntriesByEventID",
        }).Error("Failed to fetch raffle entries by event ID", err)
        return nil, err
    }
    defer rows.Close()

    var entries []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Age, &user.Gender, &user.Latitude, &user.Longitude)
        if err != nil {
            r.logger.WithFields(logrus.Fields{
                "eventId": eventID,
                "method":  "FetchRaffleEntriesByEventID",
            }).Error("Failed to scan raffle entry by event ID", err)
            return nil, err
        }
        entries = append(entries, user)
    }

    return entries, nil
}

// *************************** Helper Functions ***************************

// generateUniqueTeamID generates a unique team ID based on the current timestamp and a random number
func generateUniqueTeamID() string {
    timestamp := time.Now().Format("20060102150405")
    randomNum := rand.Intn(1000)
    return fmt.Sprintf("%s_%d", timestamp, randomNum)
}