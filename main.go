package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"event-connect/handlers"
	"event-connect/models"
	"event-connect/repositories"
	"event-connect/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize the logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Create a new router using Gorilla Mux
	r := mux.NewRouter()

	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db, logger)
	activityRepo := repositories.NewActivityRepository(db, logger)
	teamRepo := repositories.NewTeamRepository(db, logger)
	raffleRepo := repositories.NewRaffleRepository(db, logger)

	// Schedule daily team creation
    go handlers.ScheduleTeamCreation(teamRepo)

	// Initialize handlers
	eventHandler := handlers.NewEventHandler(activityRepo)

	// Middleware
	r.Use(routes.LoggingMiddleware)
	authMiddleware := routes.AuthMiddleware

	// Register routes
	routes.StaticFileRoutes(r)
	routes.HTMLFileRoutes(r)
	routes.APIRoutes(r, userRepo, activityRepo, teamRepo, raffleRepo, authMiddleware, eventHandler)
	routes.TwitterScraperRoute(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}

func initDB() (*sql.DB, error) {
	// Initialize the database connection
	db, err := models.InitializeDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
