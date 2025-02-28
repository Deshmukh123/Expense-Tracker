package models

import "time"

// Expense represents a single expense entry
type Expense struct {
    ID       int       `json:"id"`
    Amount   float64   `json:"amount"`
    Category string    `json:"category"`
    Date     time.Time `json:"date"`
}

// Category represents a spending category with a monthly budget
type Category struct {
    Name   string  `json:"name"`
    Budget float64 `json:"budget"` // Monthly budget for this category
}