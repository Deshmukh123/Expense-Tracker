package storage

import (
	"database/sql"
	"finance-tracker/models"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitializeDB connects to the MySQL database
func InitializeDB() error {
	var err error
	db, err = sql.Open("mysql", "root:root@123@tcp(127.0.0.1:3306)/finance_tracker")
	if err != nil {
		return err
	}
	return db.Ping()
}

// AddUser adds a new user to the database
func AddUser(user models.User) error {
	_, err := db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", user.Email, user.PasswordHash)
	return err
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	return user, err
}

// AddExpense adds a new expense for a user

func AddExpense(expense models.Expense, userID int) (int, error) {

	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	result, err := db.Exec("INSERT INTO expenses (amount, category, date, user_id) VALUES (?, ?, ?, ?)", expense.Amount, expense.Category, expense.Date, userID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetExpenses fetches all expenses for a user
func GetExpenses(userID int) ([]models.Expense, error) {
	rows, err := db.Query("SELECT id, amount, category, date FROM expenses WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {

		var exp models.Expense
		var dateStr string // Use string to temporarily store the date

		if err := rows.Scan(&exp.ID, &exp.Amount, &exp.Category, &dateStr); err != nil { // Scan into dateStr instead of exp.Date
			return nil, err
		}

		// Parse the date string into a time.Time value
		exp.Date, err = time.Parse("2006-01-02 15:04:05", dateStr) // Adjust the format based on your MySQL date format
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, exp)
	}
	return expenses, nil
}

// GetExpensesByCategory fetches all expenses by category for a user

func GetExpensesByCategory(category string, userID int) ([]models.Expense, error) {
	rows, err := db.Query("SELECT id, amount, category, date FROM expenses WHERE category = ? AND user_id = ?", category, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var exp models.Expense
		var dateStr string // Temporary variable to hold the date as a string

		// Scan the date column into a string
		if err := rows.Scan(&exp.ID, &exp.Amount, &exp.Category, &dateStr); err != nil {
			return nil, err
		}

		// Parse the date string into a time.Time value
		exp.Date, err = time.Parse("2006-01-02 15:04:05", dateStr) // Adjust the format based on your MySQL date format
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, exp)
	}
	return expenses, nil
}

// AddCategory adds a new category for a user
func AddCategory(category models.Category, userID int) error {
	_, err := db.Exec("INSERT INTO categories (name, budget, user_id) VALUES (?, ?, ?)", category.Name, category.Budget, userID)
	return err
}

// GetCategories fetches all categories for a user
func GetCategories(userID int) ([]models.Category, error) {
	rows, err := db.Query("SELECT id, name, budget FROM categories WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Budget); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetTotalSpendingForCategory calculates total spending for a category for a user
func GetTotalSpendingForCategory(category string, userID int) (float64, error) {
	var total float64
	err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category = ? AND user_id = ?", category, userID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
