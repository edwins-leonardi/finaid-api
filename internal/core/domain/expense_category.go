package domain

import "time"

// ExpenseCategory represents an expense category in the system
type ExpenseCategory struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateExpenseCategoryRequest represents the request to create an expense category
type CreateExpenseCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// UpdateExpenseCategoryRequest represents the request to update an expense category
type UpdateExpenseCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// ListExpenseCategoriesRequest represents the request to list expense categories
type ListExpenseCategoriesRequest struct {
	Skip  int `form:"skip"`
	Limit int `form:"limit"`
}
