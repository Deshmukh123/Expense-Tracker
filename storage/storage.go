package storage

import (
    "errors"
    "finance-tracker/models"
    "sort"
    "sync"
    "time"
)

var (
    expenses   []models.Expense
    categories []models.Category
    mu         sync.Mutex
    nextID     int = 1
)

// AddExpense adds a new expense and checks if it exceeds the category's budget
func AddExpense(expense models.Expense) (int, error) {
    mu.Lock()
    defer mu.Unlock()

    // Check if the category exists
    var category *models.Category
    for i, cat := range categories {
        if cat.Name == expense.Category {
            category = &categories[i]
            break
        }
    }
    if category == nil {
        return 0, errors.New("category does not exist")
    }

    // Check if the expense exceeds the category's budget
    totalSpent := GetTotalSpendingForCategory(category.Name)
    if totalSpent+expense.Amount > category.Budget {
        return 0, errors.New("expense exceeds category budget")
    }

    // Add the expense
    expense.ID = nextID
    nextID++
    expense.Date = time.Now()
    expenses = append(expenses, expense)
    return expense.ID, nil
}

// GetTotalSpendingForCategory calculates total spending for a category in the current month
func GetTotalSpendingForCategory(category string) float64 {
    now := time.Now()
    total := 0.0
    for _, exp := range expenses {
        if exp.Category == category && exp.Date.Month() == now.Month() && exp.Date.Year() == now.Year() {
            total += exp.Amount
        }
    }
    return total
}

// AddCategory adds a new category and sorts the list alphabetically
func AddCategory(category models.Category) {
    mu.Lock()
    defer mu.Unlock()

    // Check if the category already exists
    for _, cat := range categories {
        if cat.Name == category.Name {
            return
        }
    }

    // Add the category and sort the list
    categories = append(categories, category)
    sort.Slice(categories, func(i, j int) bool {
        return categories[i].Name < categories[j].Name
    })
}

// GetCategories returns all categories
func GetCategories() []models.Category {
    mu.Lock()
    defer mu.Unlock()
    return categories
}

// GetExpenses returns all expenses
func GetExpenses() []models.Expense {
    mu.Lock()
    defer mu.Unlock()
    return expenses
}

// GetExpensesByCategory returns all expenses for a specific category
func GetExpensesByCategory(category string) []models.Expense {
    mu.Lock()
    defer mu.Unlock()

    var filteredExpenses []models.Expense
    for _, exp := range expenses {
        if exp.Category == category {
            filteredExpenses = append(filteredExpenses, exp)
        }
    }
    return filteredExpenses
}