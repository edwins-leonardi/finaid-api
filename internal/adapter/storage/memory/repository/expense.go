package repository

import (
	"context"
	"sync"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseRepository struct {
	mu       sync.RWMutex
	expenses map[int]*domain.Expense
	nextID   int
}

// NewExpenseRepository creates a new memory expense repository
func NewExpenseRepository() port.ExpenseRepository {
	return &expenseRepository{
		expenses: make(map[int]*domain.Expense),
		nextID:   1,
	}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	expense.ID = r.nextID
	r.nextID++

	// Create a copy to avoid reference issues
	expenseCopy := *expense
	r.expenses[expense.ID] = &expenseCopy

	return nil
}

func (r *expenseRepository) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expense, exists := r.expenses[id]
	if !exists {
		return nil, domain.ErrDataNotFound
	}

	// Return a copy to avoid reference issues
	expenseCopy := *expense
	return &expenseCopy, nil
}

func (r *expenseRepository) List(ctx context.Context, filters port.ExpenseFilters) ([]*domain.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var expenses []*domain.Expense

	// Filter expenses based on criteria
	for _, expense := range r.expenses {
		if r.matchesFilters(expense, filters) {
			expenseCopy := *expense
			expenses = append(expenses, &expenseCopy)
		}
	}

	// Sort by date descending, then by created_at descending
	for i := 0; i < len(expenses)-1; i++ {
		for j := i + 1; j < len(expenses); j++ {
			if expenses[i].Date.Before(expenses[j].Date) ||
				(expenses[i].Date.Equal(expenses[j].Date) && expenses[i].CreatedAt.Before(expenses[j].CreatedAt)) {
				expenses[i], expenses[j] = expenses[j], expenses[i]
			}
		}
	}

	// Apply pagination
	start := filters.Skip
	if start >= len(expenses) {
		return []*domain.Expense{}, nil
	}

	end := start + filters.Limit
	if end > len(expenses) {
		end = len(expenses)
	}

	return expenses[start:end], nil
}

func (r *expenseRepository) matchesFilters(expense *domain.Expense, filters port.ExpenseFilters) bool {
	// Filter by category
	if filters.CategoryID != nil && expense.CategoryID != *filters.CategoryID {
		return false
	}

	// Filter by subcategory
	if filters.SubCategoryID != nil {
		if expense.SubCategoryID == nil || *expense.SubCategoryID != *filters.SubCategoryID {
			return false
		}
	}

	// Filter by payee
	if filters.PayeeID != nil && expense.PayeeID != *filters.PayeeID {
		return false
	}

	// Filter by start date
	if filters.StartDate != nil && expense.Date.Before(*filters.StartDate) {
		return false
	}

	// Filter by end date
	if filters.EndDate != nil && expense.Date.After(*filters.EndDate) {
		return false
	}

	return true
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.expenses[expense.ID]
	if !exists {
		return domain.ErrDataNotFound
	}

	// Create a copy to avoid reference issues
	expenseCopy := *expense
	r.expenses[expense.ID] = &expenseCopy

	return nil
}

func (r *expenseRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.expenses[id]
	if !exists {
		return domain.ErrDataNotFound
	}

	delete(r.expenses, id)
	return nil
}
