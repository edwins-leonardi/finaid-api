package port

import (
	"context"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

// ExpenseSubCategoryRepository defines the interface for expense subcategory data operations
type ExpenseSubCategoryRepository interface {
	Create(ctx context.Context, subcategory *domain.ExpenseSubCategory) error
	GetByID(ctx context.Context, id int) (*domain.ExpenseSubCategory, error)
	List(ctx context.Context, skip, limit int, expenseCategoryID *int) ([]*domain.ExpenseSubCategory, error)
	Update(ctx context.Context, subcategory *domain.ExpenseSubCategory) error
	Delete(ctx context.Context, id int) error
}

// ExpenseSubCategoryService defines the interface for expense subcategory business logic
type ExpenseSubCategoryService interface {
	Create(ctx context.Context, req *domain.CreateExpenseSubCategoryRequest) (*domain.ExpenseSubCategory, error)
	GetByID(ctx context.Context, id int) (*domain.ExpenseSubCategory, error)
	List(ctx context.Context, req *domain.ListExpenseSubCategoriesRequest) ([]*domain.ExpenseSubCategory, error)
	Update(ctx context.Context, id int, req *domain.UpdateExpenseSubCategoryRequest) (*domain.ExpenseSubCategory, error)
	Delete(ctx context.Context, id int) error
}
