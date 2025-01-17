package routes

import (
    "encoding/json"
    "log"
    "net/http"

    "event-connect/auth"
    "event-connect/repositories"
    "event-connect/handlers"

    "github.com/gorilla/mux"
    "github.com/justinas/alice"
)

// *************************** API Routes ***************************

// APIRoutes sets up the API routes for the application
func APIRoutes(r *mux.Router, userRepo *repositories.UserRepository, activityRepo *repositories.ActivityRepository,
    teamRepo *repositories.TeamRepository, raffleRepo *repositories.RaffleRepository, authMiddleware alice.Chain, eventHandler *handlers.EventHandler) {

    // ********** Login Route **********
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        handlers.Login(userRepo, w, r)
    }).Methods("POST")

    // ********** User Routes **********
    r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        handlers.CreateUser(userRepo)(w, r)
    }).Methods("POST")

    r.Handle("/user", authMiddleware.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.GetUserInfo(userRepo)(w, r)
    }))).Methods("GET")

    r.Handle("/profile", authMiddleware.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.GetUserProfile(userRepo, activityRepo, teamRepo)(w, r)
    }))).Methods("GET")

    r.Handle("/profile", authMiddleware.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handlers.UpdateUserProfile(userRepo)(w, r)
    }))).Methods("PUT")

    r.HandleFunc("/other-user-profile", func(w http.ResponseWriter, r *http.Request) {
        handlers.GetOtherUserProfile(userRepo)(w, r)
    }).Methods("GET")

    // ********** Event Routes **********
    r.HandleFunc("/events", eventHandler.GetEvents).Methods("GET")
    r.HandleFunc("/events/{eventId}", eventHandler.GetEventByID).Methods("GET")
    r.Handle("/events/{eventId}/register", authMiddleware.Then(http.HandlerFunc(eventHandler.RegisterEvent))).Methods("POST")

    r.HandleFunc("/events/{eventId}/user-locations", func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        eventID := params["eventId"]

        userLocations, err := activityRepo.GetUserLocationsForEvent(eventID)
        if err != nil {
            log.Printf("Error fetching user locations for event: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        log.Printf("User locations for event %s: %+v", eventID, userLocations)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userLocations)
    }).Methods("GET")

    // ********** Comment Routes **********
    r.HandleFunc("/events/{eventId}/comments", handlers.CreateComment).Methods("POST")
    r.HandleFunc("/events/{eventId}/comments", handlers.GetComments).Methods("GET")

    // ********** Team Routes **********
    r.HandleFunc("/events/{eventId}/teams", handlers.GetTeamsForEvent(teamRepo)).Methods("GET")
    r.HandleFunc("/trigger-create-teams/{eventId}", handlers.TriggerCreateTeams(teamRepo)).Methods("POST")

    // ********** Raffle Routes **********
    r.HandleFunc("/events/{eventId}/raffle", func(w http.ResponseWriter, r *http.Request) {
        handlers.EnterRaffle(raffleRepo, auth.GetUserIDFromToken)(w, r)
    }).Methods("POST")
}