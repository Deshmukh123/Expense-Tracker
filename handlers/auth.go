package handlers

import (
	"encoding/json"
	"finance-tracker/models"
	"finance-tracker/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func Register(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body into the credentials struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Registering user with email: %s", credentials.Email)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create a new user with the hashed password
	user := models.User{
		Email:        credentials.Email,
		PasswordHash: string(hashedPassword),
	}

	log.Printf("Hashed password: %s", user.PasswordHash)

	// Save the user to the database
	if err := storage.AddUser(user); err != nil {
		log.Printf("Failed to add user to database: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User registered successfully: %s", user.Email)
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body into the credentials struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for email: %s", credentials.Email)

	// Fetch the user from the database
	user, err := storage.GetUserByEmail(credentials.Email)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	log.Printf("User found: %s", user.Email)
	log.Printf("Stored hashed password: %s", user.PasswordHash)
	log.Printf("Provided password: %s", credentials.Password)

	// Compare the password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		log.Printf("Password mismatch for user: %s", user.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	log.Printf("Password verified for user: %s", user.Email)

	// Generate a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("Token generated for user: %s", user.Email)

	// Return the token
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
