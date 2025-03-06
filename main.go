package main

import (
	"finance-tracker/handlers"
	"finance-tracker/middleware"
	"finance-tracker/storage"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var jwtKey []byte

func main() {

	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the JWT key from the environment variable
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		log.Fatal("JWT_SECRET_KEY is not set. Please set it in the .env file.")
	}
	log.Printf("JWT key: %s\n", jwtKey)


	// Initialize the database
	if err := storage.InitializeDB(); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/reports/monthly", handlers.HandleMonthlyReport)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("/expenses", handlers.HandleExpenses)
	protected.HandleFunc("/categories", handlers.HandleCategories)
	protected.HandleFunc("/spending/category", handlers.HandleSpendingByCategory)

	http.Handle("/", middleware.AuthMiddleware(protected))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
