package handlers

import (
	"encoding/json"
	"finance-tracker/models"
	"finance-tracker/storage"
	"log"
	"net/http"
)

// HandleExpenses handles requests for expenses
func HandleExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	switch r.Method {
	case http.MethodPost:
		var expense models.Expense
		if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := storage.AddExpense(expense, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	case http.MethodGet:
		expenses, err := storage.GetExpenses(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expenses)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleCategories handles requests for categories
func HandleCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	switch r.Method {
	case http.MethodPost:
		var category models.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := storage.AddCategory(category, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		categories, err := storage.GetCategories(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleSpendingByCategory handles requests for spending by category
func HandleSpendingByCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "category parameter is required", http.StatusBadRequest)
		return
	}

	totalSpent, err := storage.GetTotalSpendingForCategory(category, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"total_spent": totalSpent})
}

//Handles monthly report

func HandleMonthlyReport(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Println("user_id is nil in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		log.Println("user_id is not an int")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if month == "" || year == "" {
		http.Error(w, "month and year parameters are required", http.StatusBadRequest)
		return
	}

	report, err := storage.GetMonthlyReport(userIDInt, month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
