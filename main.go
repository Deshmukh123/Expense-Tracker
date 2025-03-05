package main

import (
    "finance-tracker/handlers"
    "finance-tracker/middleware"
    "finance-tracker/storage"
    "log"
    "net/http"
)

func main() {
    // Initialize the database
    if err := storage.InitializeDB(); err != nil {
        log.Fatalf("Could not initialize database: %v", err)
    }

    // Set up HTTP routes
    http.HandleFunc("/register", handlers.Register)
    http.HandleFunc("/login", handlers.Login)

    // Protected routes
    protected := http.NewServeMux()
    protected.HandleFunc("/expenses", handlers.HandleExpenses)
    protected.HandleFunc("/categories", handlers.HandleCategories)
    protected.HandleFunc("/spending/category", handlers.HandleSpendingByCategory)

    http.Handle("/", middleware.AuthMiddleware(protected))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}