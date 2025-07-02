package domain

import "time"

// Expense represents an expense in the system
type Expense struct {
	ID            int       `json:"id"`
	Amount        float64   `json:"amount"`
	CategoryID    int       `json:"category_id"`
	SubCategoryID *int      `json:"subcategory_id,omitempty"` // Optional
	Date          time.Time `json:"date"`
	PayeeID       int       `json:"payee_id"`   // Person who received the payment
	AccountID     int       `json:"account_id"` // Account from which the expense was paid
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateExpenseRequest represents the request to create an expense
type CreateExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=0"`
	CategoryID    int     `json:"category_id" binding:"required,min=1"`
	SubCategoryID *int    `json:"subcategory_id,omitempty"`
	Date          string  `json:"date" binding:"required"` // Format: YYYY-MM-DD
	PayeeID       int     `json:"payee_id" binding:"required,min=1"`
	AccountID     int     `json:"account_id" binding:"required,min=1"`
	Notes         string  `json:"notes,omitempty"`
}

// UpdateExpenseRequest represents the request to update an expense
type UpdateExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=0"`
	CategoryID    int     `json:"category_id" binding:"required,min=1"`
	SubCategoryID *int    `json:"subcategory_id,omitempty"`
	Date          string  `json:"date" binding:"required"` // Format: YYYY-MM-DD
	PayeeID       int     `json:"payee_id" binding:"required,min=1"`
	AccountID     int     `json:"account_id" binding:"required,min=1"`
	Notes         string  `json:"notes,omitempty"`
}

// ListExpensesRequest represents the request to list expenses
type ListExpensesRequest struct {
	Skip          int    `form:"skip"`
	Limit         int    `form:"limit"`
	CategoryID    int    `form:"category_id"`    // Optional filter by category
	SubCategoryID int    `form:"subcategory_id"` // Optional filter by subcategory
	PayeeID       int    `form:"payee_id"`       // Optional filter by payee
	AccountID     int    `form:"account_id"`     // Optional filter by account
	StartDate     string `form:"start_date"`     // Optional filter by date range (YYYY-MM-DD)
	EndDate       string `form:"end_date"`       // Optional filter by date range (YYYY-MM-DD)
}
