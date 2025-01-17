// handlers/auth.go
package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("QXR0cUZYZjM5STB3NzNubE1jdVhYTVBxNjhQR3JaWWF5NGt4RkZYRXlwcHhTSEQ2dw==")

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"userId"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint) (string, error) {
	claims := &Claims{
		UserID: strconv.FormatUint(uint64(userID), 10),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
		return "", err
	}

	log.Printf("Generated token: %s\n", signedToken)
	log.Printf("Secret key used for signing: %s\n", string(jwtKey))

	return signedToken, nil
}

// GetUserIDFromToken retrieves the user ID from the provided HTTP request's Authorization header token.
func GetUserIDFromToken(r *http.Request) (int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Println("Authorization header missing")
		return 0, errors.New("authorization header missing")
	}

	// Remove the "Bearer " prefix from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v\n", err)
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Extract the user ID from the token claims
		userID, err := strconv.Atoi(claims.UserID)
		if err != nil {
			log.Printf("Error converting user ID to integer: %v\n", err)
			return 0, err
		}
		return userID, nil
	}
	log.Println("Invalid token")
	return 0, errors.New("invalid token")
}

// VerifyToken verifies the provided JWT token and returns the user ID if the token is valid.
func VerifyToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID, err := strconv.Atoi(claims.UserID)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
