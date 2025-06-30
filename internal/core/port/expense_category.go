package port

import (
	"context"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

// ExpenseCategoryRepository defines the interface for expense category data operations
type ExpenseCategoryRepository interface {
	Create(ctx context.Context, category *domain.ExpenseCategory) error
	GetByID(ctx context.Context, id int) (*domain.ExpenseCategory, error)
	List(ctx context.Context, skip, limit int) ([]*domain.ExpenseCategory, error)
	Update(ctx context.Context, category *domain.ExpenseCategory) error
	Delete(ctx context.Context, id int) error
}

// ExpenseCategoryService defines the interface for expense category business logic
type ExpenseCategoryService interface {
	Create(ctx context.Context, req *domain.CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error)
	GetByID(ctx context.Context, id int) (*domain.ExpenseCategory, error)
	List(ctx context.Context, req *domain.ListExpenseCategoriesRequest) ([]*domain.ExpenseCategory, error)
	Update(ctx context.Context, id int, req *domain.UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error)
	Delete(ctx context.Context, id int) error
}
