package handlers

import (
	"encoding/json"
	"event-connect/auth"
	"event-connect/models"
	"event-connect/repositories"
	"event-connect/skiddle"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// *************************** EventHandler ***************************

// EventHandler represents the handler for event-related operations
type EventHandler struct {
	activityRepo *repositories.ActivityRepository
}

// NewEventHandler creates a new instance of EventHandler
func NewEventHandler(activityRepo *repositories.ActivityRepository) *EventHandler {
	return &EventHandler{
		activityRepo: activityRepo,
	}
}

// *************************** Handler Methods ***************************

// GetEventByID retrieves event details by event ID
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID, err := strconv.ParseUint(params["eventId"], 10, 64)
	if err != nil {
		log.Printf("Error parsing event ID: %v", err)
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	url := skiddle.EventDetailsURL(uint(eventID))
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching event details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error parsing event details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resultsObj, ok := result["results"].(map[string]interface{})
	if !ok {
		log.Println("Invalid response format: missing 'results' object")
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	eventData := map[string]interface{}{
		"id":          resultsObj["id"],
		"eventname":   resultsObj["eventname"],
		"date":        resultsObj["date"],
		"venue":       resultsObj["venue"],
		"description": resultsObj["description"],
		"entryprice":  resultsObj["entryprice"],
		"minage":      resultsObj["MinAge"],
		"link":        resultsObj["link"],
		"imageurl":    resultsObj["imageurl"],
	}

	log.Printf("Event details: %+v", eventData)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(eventData); err != nil {
		log.Printf("Error encoding event details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetEvents retrieves events based on query parameters
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	params := make(map[string]string)
	params["order"] = "date"
	params["description"] = "1"
	if latitude := queryParams.Get("latitude"); latitude != "" {
		params["latitude"] = latitude
	}
	if longitude := queryParams.Get("longitude"); longitude != "" {
		params["longitude"] = longitude
	}
	if radius := queryParams.Get("radius"); radius != "" {
		params["radius"] = radius
	}
	if eventcode := queryParams.Get("eventcode"); eventcode != "" {
		params["eventcode"] = eventcode
	}
	if keyword := queryParams.Get("keyword"); keyword != "" {
		params["keyword"] = keyword
	}
	if minDate := queryParams.Get("minDate"); minDate != "" {
		params["minDate"] = minDate
	}
	if maxDate := queryParams.Get("maxDate"); maxDate != "" {
		params["maxDate"] = maxDate
	}

	url := skiddle.EventSearchURL(params)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching events: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error parsing events: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	results, ok := result["results"].([]interface{})
	if !ok {
		log.Println("Invalid response format: missing 'results' array")
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	var events []map[string]interface{}
	for _, item := range results {
		event, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		venue, ok := event["venue"].(map[string]interface{})
		if !ok {
			continue
		}
		latitude, _ := venue["latitude"].(float64)
		longitude, _ := venue["longitude"].(float64)

		weather, err := models.GetWeather(fmt.Sprintf("%.6f", latitude), fmt.Sprintf("%.6f", longitude))
		if err != nil {
			log.Printf("Error fetching weather data for event %v: %v", event["id"], err)
			continue
		}

		eventData := map[string]interface{}{
			"id":          event["id"],
			"eventname":   event["eventname"],
			"date":        event["date"],
			"venue":       event["venue"],
			"description": event["description"],
			"entryprice":  event["entryprice"],
			"minage":      event["minage"],
			"link":        event["link"],
			"imageurl":    event["imageurl"],
			"weather":     weather.Weather,
			"temperature": weather.Temperature,
		}
		events = append(events, eventData)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("Error encoding events: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RegisterEvent registers a user for an event
func (h *EventHandler) RegisterEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	activity := &models.Activity{
		UserID:       uint(userID),
		EventID:      uint(eventID),
		ActivityType: "event_registered",
		Timestamp:    time.Now(),
	}
	err = h.activityRepo.CreateActivity(activity)
	if err != nil {
		log.Printf("Error creating activity record: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}