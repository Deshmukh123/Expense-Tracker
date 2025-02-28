package main

import (
    "finance-tracker/handlers"
    "log"
    "net/http"
)

func main() {
    // Endpoint to add and list expenses
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

    // Endpoint to get expenses by category
    http.HandleFunc("/expenses/category", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handlers.GetExpensesByCategory(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Endpoint to add and list categories
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

    // Endpoint to get total spending for a category
    http.HandleFunc("/spending/category", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handlers.GetTotalSpendingForCategory(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Start the server
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}