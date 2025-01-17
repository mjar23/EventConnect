// models/db_initializer.go
package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func GetDatabaseConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "admin"
	}
	if dbname == "" {
		dbname = "postgres"
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return connectionString
}

func InitializeDB() (*sql.DB, error) {
	connectionString := GetDatabaseConnectionString()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Create users table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL UNIQUE,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
            first_name VARCHAR(255),
            last_name VARCHAR(255),
            bio TEXT,
            interests TEXT,
            location TEXT,
            latitude DOUBLE PRECISION,
            longitude DOUBLE PRECISION,
            age INTEGER,
            gender TEXT,
            age_min INTEGER,
            age_max INTEGER,
            distance_preference INTEGER,
            instagram_username VARCHAR(255),
            facebook_username VARCHAR(255),
            snapchat_username VARCHAR(255)
        )
    `)
	if err != nil {
		return nil, err
	}

	// Create activities table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS activities (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            event_id INTEGER,
            activity_type VARCHAR(50),
            timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, err
	}

	// Create comments table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            id SERIAL PRIMARY KEY,
            event_id VARCHAR(255),
            user_id INTEGER REFERENCES users(id),
            text TEXT NOT NULL,
            created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return nil, err
	}

	// Create raffle_entries table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS raffle_entries (
            id SERIAL PRIMARY KEY,
            event_id VARCHAR(255) NOT NULL,
            user_id INTEGER NOT NULL REFERENCES users(id),
            age INTEGER NOT NULL,
            gender VARCHAR(10) NOT NULL,
            latitude DOUBLE PRECISION NOT NULL,
            longitude DOUBLE PRECISION NOT NULL
        )
    `)
	if err != nil {
		return nil, err
	}

	// Create teams table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS teams (
        id SERIAL PRIMARY KEY,
        event_id INTEGER,
        user_id INTEGER REFERENCES users(id),
        age INTEGER,
        gender VARCHAR(10),
        latitude DOUBLE PRECISION,
        longitude DOUBLE PRECISION,
        created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        team_id VARCHAR(20) NOT NULL
    )
`)
	if err != nil {
		return nil, err
	}

	log.Println("Database tables initialized successfully")
	return db, nil
}
