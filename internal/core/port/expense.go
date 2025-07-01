package port

import (
	"context"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

// ExpenseRepository defines the interface for expense data operations
type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	GetByID(ctx context.Context, id int) (*domain.Expense, error)
	List(ctx context.Context, filters ExpenseFilters) ([]*domain.Expense, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id int) error
}

// ExpenseFilters represents filters for listing expenses
type ExpenseFilters struct {
	Skip          int
	Limit         int
	CategoryID    *int
	SubCategoryID *int
	PayeeID       *int
	StartDate     *time.Time
	EndDate       *time.Time
}

// ExpenseService defines the interface for expense business logic
type ExpenseService interface {
	Create(ctx context.Context, req *domain.CreateExpenseRequest) (*domain.Expense, error)
	GetByID(ctx context.Context, id int) (*domain.Expense, error)
	List(ctx context.Context, req *domain.ListExpensesRequest) ([]*domain.Expense, error)
	Update(ctx context.Context, id int, req *domain.UpdateExpenseRequest) (*domain.Expense, error)
	Delete(ctx context.Context, id int) error
}
