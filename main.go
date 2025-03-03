package main

import (
	"finance-tracker/handlers"
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
	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddExpense(w, r)
		case http.MethodGet:
			handlers.GetExpenses(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddCategory(w, r)
		case http.MethodGet:
			handlers.GetCategories(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/spending/category", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTotalSpendingForCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}