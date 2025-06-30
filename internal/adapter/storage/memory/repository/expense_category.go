package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseCategoryRepository struct {
	categories map[int]*domain.ExpenseCategory
	nextID     int
	mu         sync.RWMutex
}

// NewExpenseCategoryRepository creates a new in-memory expense category repository
func NewExpenseCategoryRepository() port.ExpenseCategoryRepository {
	return &expenseCategoryRepository{
		categories: make(map[int]*domain.ExpenseCategory),
		nextID:     1,
	}
}

func (r *expenseCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	category.ID = r.nextID
	r.nextID++

	// Create a copy to avoid reference issues
	categoryCopy := &domain.ExpenseCategory{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	r.categories[category.ID] = categoryCopy
	return nil
}

func (r *expenseCategoryRepository) GetByID(ctx context.Context, id int) (*domain.ExpenseCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, exists := r.categories[id]
	if !exists {
		return nil, domain.ErrDataNotFound
	}

	// Return a copy to avoid reference issues
	return &domain.ExpenseCategory{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (r *expenseCategoryRepository) List(ctx context.Context, skip, limit int) ([]*domain.ExpenseCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	categories := make([]*domain.ExpenseCategory, 0, len(r.categories))
	for _, category := range r.categories {
		categories = append(categories, &domain.ExpenseCategory{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
	}

	// Sort by created_at descending (newest first)
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].CreatedAt.After(categories[j].CreatedAt)
	})

	// Apply pagination
	start := skip
	if start > len(categories) {
		start = len(categories)
	}

	end := start + limit
	if end > len(categories) {
		end = len(categories)
	}

	return categories[start:end], nil
}

func (r *expenseCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.categories[category.ID]
	if !exists {
		return domain.ErrDataNotFound
	}

	// Update the category
	r.categories[category.ID] = &domain.ExpenseCategory{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: time.Now(),
	}

	return nil
}

func (r *expenseCategoryRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.categories[id]
	if !exists {
		return domain.ErrDataNotFound
	}

	delete(r.categories, id)
	return nil
}
