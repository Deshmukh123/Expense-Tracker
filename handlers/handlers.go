package handlers

import (
    "encoding/json"
    "finance-tracker/models"
    "finance-tracker/storage"
    "net/http"
)

// AddExpense handles adding a new expense
func AddExpense(w http.ResponseWriter, r *http.Request) {
    var expense models.Expense
    if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id, err := storage.AddExpense(expense)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetExpenses returns all expenses
func GetExpenses(w http.ResponseWriter, r *http.Request) {
    expenses := storage.GetExpenses()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(expenses)
}

// GetExpensesByCategory returns all expenses for a specific category
func GetExpensesByCategory(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category")
    if category == "" {
        http.Error(w, "category parameter is required", http.StatusBadRequest)
        return
    }

    expenses := storage.GetExpensesByCategory(category)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(expenses)
}

// AddCategory handles adding a new category
func AddCategory(w http.ResponseWriter, r *http.Request) {
    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    storage.AddCategory(category)
    w.WriteHeader(http.StatusCreated)
}

// GetCategories returns all categories
func GetCategories(w http.ResponseWriter, r *http.Request) {
    categories := storage.GetCategories()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

// GetTotalSpendingForCategory returns the total spending for a category in the current month
func GetTotalSpendingForCategory(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category")
    if category == "" {
        http.Error(w, "category parameter is required", http.StatusBadRequest)
        return
    }

    totalSpent := storage.GetTotalSpendingForCategory(category)
    json.NewEncoder(w).Encode(map[string]float64{"total_spent": totalSpent})
}