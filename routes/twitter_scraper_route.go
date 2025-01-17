package routes

import (
    "encoding/json"
    "log"
    "net/http"

    "event-connect/twitter"

    "github.com/gorilla/mux"
)

func TwitterScraperRoute(r *mux.Router) {
    // Add a route for triggering Twitter scraping
    r.HandleFunc("/events/{eventName}/twitter-scraper", func(w http.ResponseWriter, r *http.Request) {
        // Retrieve the event name from URL parameters
        eventName := mux.Vars(r)["eventName"]

        log.Printf("Scraping tweets for event: %s", eventName)

        // Call the Twitter scraper function passing the event name
        tweets, err := twitter.Twitterscrapering(eventName)
        if err != nil {
            log.Printf("Error scraping tweets: %v", err)
            http.Error(w, "Failed to scrape tweets", http.StatusInternalServerError)
            return
        }

        // Return the scraped tweets as JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(tweets); err != nil {
            log.Printf("Error encoding tweets to JSON: %v", err)
            http.Error(w, "Failed to encode tweets", http.StatusInternalServerError)
        }
    }).Methods("GET")
}