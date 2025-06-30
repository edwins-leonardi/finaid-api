package domain

import "time"

// ExpenseSubCategory represents an expense subcategory in the system
type ExpenseSubCategory struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	ExpenseCategoryID int       `json:"expense_category_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateExpenseSubCategoryRequest represents the request to create an expense subcategory
type CreateExpenseSubCategoryRequest struct {
	Name              string `json:"name" binding:"required,min=1,max=100"`
	ExpenseCategoryID int    `json:"expense_category_id" binding:"required,min=1"`
}

// UpdateExpenseSubCategoryRequest represents the request to update an expense subcategory
type UpdateExpenseSubCategoryRequest struct {
	Name              string `json:"name" binding:"required,min=1,max=100"`
	ExpenseCategoryID int    `json:"expense_category_id" binding:"required,min=1"`
}

// ListExpenseSubCategoriesRequest represents the request to list expense subcategories
type ListExpenseSubCategoriesRequest struct {
	Skip              int `form:"skip"`
	Limit             int `form:"limit"`
	ExpenseCategoryID int `form:"expense_category_id"` // Optional filter by category
}
