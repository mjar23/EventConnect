package handlers

import (
	"encoding/json"
	"event-connect/models"
	
	"event-connect/repositories"
	"event-connect/emailUtil"
	"event-connect/skiddle"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func haversineDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const earthRadius = 6371 // Radius of the Earth in kilometers

	lat1 = degToRad(lat1)
	lon1 = degToRad(lon1)
	lat2 = degToRad(lat2)
	lon2 = degToRad(lon2)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func groupEntries(entries []models.User) []models.Team {
	var teams []models.Team

	// Group entries by gender
	entriesByGender := make(map[string][]models.User)
	for _, entry := range entries {
		entriesByGender[entry.Gender] = append(entriesByGender[entry.Gender], entry)
	}

	// Create teams for each gender
	for _, genderEntries := range entriesByGender {
		// Sort entries by age and location (latitude, longitude)
		sort.Slice(genderEntries, func(i, j int) bool {
			if genderEntries[i].Age == genderEntries[j].Age {
				distanceI := haversineDistance(genderEntries[i].Longitude, genderEntries[i].Latitude, 0, 0)
				distanceJ := haversineDistance(genderEntries[j].Longitude, genderEntries[j].Latitude, 0, 0)
				return distanceI < distanceJ
			}
			return genderEntries[i].Age < genderEntries[j].Age
		})

		var currentTeam []models.Member
		for _, entry := range genderEntries {
			if len(currentTeam) < 4 {
				currentTeam = append(currentTeam, models.Member{
					UserID:    entry.ID,
					Username:  entry.Username,
					Age:       entry.Age,
					Gender:    entry.Gender,
					Latitude:  entry.Latitude,
					Longitude: entry.Longitude,
				})
			}

			if len(currentTeam) == 4 {
				teams = append(teams, models.Team{Members: currentTeam})
				currentTeam = nil
			}
		}

		if len(currentTeam) > 0 {
			teams = append(teams, models.Team{Members: currentTeam})
		}
	}

	return teams
}

func CreateTeams(teamRepo *repositories.TeamRepository, eventID uint) error {
	log.Printf("Creating teams for event ID: %d", eventID)

	// Fetch raffle entries for the given event ID
	entries, err := teamRepo.FetchRaffleEntries(eventID)
	if err != nil {
		return fmt.Errorf("failed to fetch raffle entries: %w", err)
	}
	log.Printf("Fetched %d raffle entries for event ID: %d", len(entries), eventID)

	// Group entries into teams
	teams := groupEntries(entries)
	log.Printf("Created %d teams for event ID: %d", len(teams), eventID)

	// Insert teams into the database
	err = teamRepo.InsertTeams(eventID, teams)
	if err != nil {
		return fmt.Errorf("failed to insert teams: %w", err)
	}
	log.Printf("Inserted teams into the database for event ID: %d", eventID)

	// Send emails to team members
	err = sendTeamEmails(teamRepo, teams)
	if err != nil {
		return fmt.Errorf("failed to send team emails: %w", err)
	}
	log.Printf("Sent team emails for event ID: %d", eventID)

	return nil
}

func GetUserTeams(teamRepo *repositories.TeamRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userId")
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		teams, err := teamRepo.FetchUserTeams(uint(userID))
		if err != nil {
			http.Error(w, "Failed to fetch user teams", http.StatusInternalServerError)
			return
		}

		for i := range teams {
			eventDetails, err := fetchEventDetails(teams[i].EventID)
			if err != nil {
				log.Printf("Failed to fetch event details for event ID %d: %v", teams[i].EventID, err)
				continue
			}
			teams[i].EventName = eventDetails["eventname"].(string)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teams)
	}
}

func GetTeamsForEvent(teamRepo *repositories.TeamRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		eventIDStr := params["eventId"]
		eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		teams, err := teamRepo.GetTeamsForEvent(uint(eventID))
		if err != nil {
			http.Error(w, "Failed to fetch teams for event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teams)
	}
}

func TriggerCreateTeams(teamRepo *repositories.TeamRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		eventIDStr := params["eventId"]
		eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		err = CreateTeams(teamRepo, uint(eventID))
		if err != nil {
			http.Error(w, "Failed to create teams", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Teams created successfully for event ID: %s", eventIDStr)
	}
}

func ScheduleTeamCreation(teamRepo *repositories.TeamRepository) {
	ticker := time.NewTicker(24 * time.Hour) // Run the task every 24 hours
	defer ticker.Stop()

	for range ticker.C {
		createTeamsForUpcomingEvents(teamRepo)
	}
}

func createTeamsForUpcomingEvents(teamRepo *repositories.TeamRepository) {
	log.Printf("Checking raffle entries...")

	eventIDs, err := teamRepo.FetchEventIDsFromRaffleEntries()
	if err != nil {
		log.Printf("Error fetching event IDs: %v", err)
		return
	}
	for _, eventID := range eventIDs {
		log.Printf("Checking event ID: %d", eventID)

		// Check if the event is exactly 1 week away using the Skiddle API
		eventDetails, err := fetchEventDetails(eventID)
		if err != nil {
			log.Printf("Error checking event in Skiddle API: %v", err)
			continue
		}

		eventDateStr, ok := eventDetails["date"].(string)
		if !ok {
			log.Printf("Event date not found in Skiddle API response for event ID %d", eventID)
			continue
		}

		eventDate, err := time.Parse("2006-01-02", eventDateStr)
		if err != nil {
			log.Printf("Error parsing event date: %v", err)
			continue
		}

		oneWeekFromToday := time.Now().AddDate(0, 0, 7)
		eventDiff := eventDate.Sub(oneWeekFromToday)
		if eventDiff >= -24*time.Hour && eventDiff < 24*time.Hour {
			log.Printf("Event ID %d is exactly 1 week away", eventID)

			// Fetch users related to the event ID from the raffle_entries table
			entries, err := teamRepo.FetchRaffleEntriesByEventID(eventID)
			if err != nil {
				log.Printf("Error fetching raffle entries for event ID %d: %v", eventID, err)
				continue
			}

			log.Printf("Fetched %d raffle entries for event ID %d", len(entries), eventID)

			// Create teams using the fetched users
			teams := groupEntries(entries)
			log.Printf("Created %d teams for event ID %d", len(teams), eventID)

			// Insert teams into the database
			err = teamRepo.InsertTeams(eventID, teams)
			if err != nil {
				log.Printf("Error inserting teams for event ID %d: %v", eventID, err)
				continue
			}

			log.Printf("Teams created successfully for event ID %d", eventID)
		} else {
			log.Printf("Event ID %d is not exactly 1 week away", eventID)
		}
	}
}
func fetchEventDetails(eventID uint) (map[string]interface{}, error) {
    url := skiddle.EventDetailsURL(eventID)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    eventDetails, ok := result["results"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid event details response")
    }

    return eventDetails, nil
}
func sendTeamEmails(teamRepo *repositories.TeamRepository, teams []models.Team) error {
	for _, team := range teams {
		var teamMemberEmails []string
		var teamMemberSocials []string
		for _, member := range team.Members {
			user, err := teamRepo.GetUserByID(member.UserID)
			if err != nil {
				return err
			}

			teamMemberEmails = append(teamMemberEmails, user.Email)
			teamMemberSocials = append(teamMemberSocials, "Instagram: "+user.InstagramUsername+", Facebook: "+user.FacebookUsername+", Snapchat: "+user.SnapchatUsername)
		}

		// Construct the email message
		emailSubject := "Your Team for Event ID: " + strconv.Itoa(int(team.EventID))
		emailBody := "Dear team members,\n\nYour team for the event has been created. The members of your team are:\n\n"
		for _, member := range team.Members {
			emailBody += "- " + member.Username + " (Age: " + strconv.Itoa(member.Age) + ", Gender: " + member.Gender + ")\n"
		}
		emailBody += "\nYour team's social media usernames are:\n\n"
		for _, social := range teamMemberSocials {
			emailBody += "- " + social + "\n"
		}
		emailBody += "\nBest regards,\nThe Event Team"

		// Send the email
		err := sendEmail(teamMemberEmails, emailSubject, emailBody)
		if err != nil {
			return err
		}
	}

	return nil
}
func sendEmail(recipients []string, subject, body string) error {
    // Create HTML content (you can make this more sophisticated if needed)
    htmlContent := fmt.Sprintf("<html><body><p>%s</p></body></html>", body)

    // Use the SendEmail function from emailUtil
    err := emailUtil.SendEmail(recipients, subject, body, htmlContent)
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}
