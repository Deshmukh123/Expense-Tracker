package models

import "time"

type User struct {
    ID           int    `json:"id"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"`
}

type Expense struct {
    ID       int       `json:"id"`
    Amount   float64   `json:"amount"`
    Category string    `json:"category"`
    Date     time.Time `json:"date"`
    UserID   int       `json:"user_id"`
}

type Category struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Budget float64 `json:"budget"`
    UserID int    `json:"user_id"`
}