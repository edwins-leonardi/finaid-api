package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseSubCategoryRepository struct {
	subcategories map[int]*domain.ExpenseSubCategory
	nextID        int
	mu            sync.RWMutex
}

// NewExpenseSubCategoryRepository creates a new in-memory expense subcategory repository
func NewExpenseSubCategoryRepository() port.ExpenseSubCategoryRepository {
	return &expenseSubCategoryRepository{
		subcategories: make(map[int]*domain.ExpenseSubCategory),
		nextID:        1,
	}
}

func (r *expenseSubCategoryRepository) Create(ctx context.Context, subcategory *domain.ExpenseSubCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	subcategory.ID = r.nextID
	r.nextID++

	// Create a copy to avoid reference issues
	subcategoryCopy := &domain.ExpenseSubCategory{
		ID:                subcategory.ID,
		Name:              subcategory.Name,
		ExpenseCategoryID: subcategory.ExpenseCategoryID,
		CreatedAt:         subcategory.CreatedAt,
		UpdatedAt:         subcategory.UpdatedAt,
	}

	r.subcategories[subcategory.ID] = subcategoryCopy
	return nil
}

func (r *expenseSubCategoryRepository) GetByID(ctx context.Context, id int) (*domain.ExpenseSubCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subcategory, exists := r.subcategories[id]
	if !exists {
		return nil, domain.ErrDataNotFound
	}

	// Return a copy to avoid reference issues
	return &domain.ExpenseSubCategory{
		ID:                subcategory.ID,
		Name:              subcategory.Name,
		ExpenseCategoryID: subcategory.ExpenseCategoryID,
		CreatedAt:         subcategory.CreatedAt,
		UpdatedAt:         subcategory.UpdatedAt,
	}, nil
}

func (r *expenseSubCategoryRepository) List(ctx context.Context, skip, limit int, expenseCategoryID *int) ([]*domain.ExpenseSubCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice and apply filter if needed
	subcategories := make([]*domain.ExpenseSubCategory, 0, len(r.subcategories))
	for _, subcategory := range r.subcategories {
		// Apply filter by expense category if specified
		if expenseCategoryID != nil && subcategory.ExpenseCategoryID != *expenseCategoryID {
			continue
		}

		subcategories = append(subcategories, &domain.ExpenseSubCategory{
			ID:                subcategory.ID,
			Name:              subcategory.Name,
			ExpenseCategoryID: subcategory.ExpenseCategoryID,
			CreatedAt:         subcategory.CreatedAt,
			UpdatedAt:         subcategory.UpdatedAt,
		})
	}

	// Sort by created_at descending (newest first)
	sort.Slice(subcategories, func(i, j int) bool {
		return subcategories[i].CreatedAt.After(subcategories[j].CreatedAt)
	})

	// Apply pagination
	start := skip
	if start > len(subcategories) {
		start = len(subcategories)
	}

	end := start + limit
	if end > len(subcategories) {
		end = len(subcategories)
	}

	return subcategories[start:end], nil
}

func (r *expenseSubCategoryRepository) Update(ctx context.Context, subcategory *domain.ExpenseSubCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.subcategories[subcategory.ID]
	if !exists {
		return domain.ErrDataNotFound
	}

	// Update the subcategory
	r.subcategories[subcategory.ID] = &domain.ExpenseSubCategory{
		ID:                subcategory.ID,
		Name:              subcategory.Name,
		ExpenseCategoryID: subcategory.ExpenseCategoryID,
		CreatedAt:         subcategory.CreatedAt,
		UpdatedAt:         time.Now(),
	}

	return nil
}

func (r *expenseSubCategoryRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.subcategories[id]
	if !exists {
		return domain.ErrDataNotFound
	}

	delete(r.subcategories, id)
	return nil
}
