package storage

import (
	"database/sql"
	"errors"
	"finance-tracker/models"

	// "fmt"
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

// AddExpense adds a new expense to the database
func AddExpense(expense models.Expense) (int, error) {
	// Check if the category exists
	var budget float64
	err := db.QueryRow("SELECT budget FROM categories WHERE name = ?", expense.Category).Scan(&budget)
	if err != nil {
		return 0, errors.New("category does not exist")
	}

	// Check if the expense exceeds the category's budget
	totalSpent, err := GetTotalSpendingForCategory(expense.Category)
	if err != nil {
		return 0, err
	}
	if totalSpent+expense.Amount > budget {
		return 0, errors.New("expense exceeds category budget")
	}

	// Insert the expense
	result, err := db.Exec("INSERT INTO expenses (amount, category, date) VALUES (?, ?, ?)", expense.Amount, expense.Category, time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetTotalSpendingForCategory calculates total spending for a category in the current month
func GetTotalSpendingForCategory(category string) (float64, error) {
	var total float64
	err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category = ? AND MONTH(date) = MONTH(CURRENT_DATE()) AND YEAR(date) = YEAR(CURRENT_DATE())", category).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// AddCategory adds a new category to the database
func AddCategory(category models.Category) error {
	_, err := db.Exec("INSERT INTO categories (name, budget) VALUES (?, ?)", category.Name, category.Budget)
	return err
}

// GetCategories returns all categories from the database
func GetCategories() ([]models.Category, error) {
	rows, err := db.Query("SELECT name, budget FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.Name, &cat.Budget); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetExpenses returns all expenses from the database
func GetExpenses() ([]models.Expense, error) {
	rows, err := db.Query("SELECT id, amount, category, date FROM expenses")
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

// GetExpensesByCategory returns all expenses for a specific category
func GetExpensesByCategory(category string) ([]models.Expense, error) {
	rows, err := db.Query("SELECT id, amount, category, date FROM expenses WHERE category = ?", category)
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
