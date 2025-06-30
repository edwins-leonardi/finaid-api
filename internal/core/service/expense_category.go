package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseCategoryService struct {
	repo   port.ExpenseCategoryRepository
	logger *slog.Logger
}

// NewExpenseCategoryService creates a new expense category service
func NewExpenseCategoryService(repo port.ExpenseCategoryRepository, logger *slog.Logger) port.ExpenseCategoryService {
	return &expenseCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *expenseCategoryService) Create(ctx context.Context, req *domain.CreateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	s.logger.Info("Creating expense category", "name", req.Name)

	// Validate and sanitize input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	category := &domain.ExpenseCategory{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		s.logger.Error("Failed to create expense category", "error", err, "name", name)
		return nil, err
	}

	s.logger.Info("Expense category created successfully", "id", category.ID, "name", category.Name)
	return category, nil
}

func (s *expenseCategoryService) GetByID(ctx context.Context, id int) (*domain.ExpenseCategory, error) {
	s.logger.Info("Getting expense category by ID", "id", id)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense category", "error", err, "id", id)
		return nil, err
	}

	return category, nil
}

func (s *expenseCategoryService) List(ctx context.Context, req *domain.ListExpenseCategoriesRequest) ([]*domain.ExpenseCategory, error) {
	s.logger.Info("Listing expense categories", "skip", req.Skip, "limit", req.Limit)

	// Set default values
	skip := req.Skip
	if skip < 0 {
		skip = 0
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	categories, err := s.repo.List(ctx, skip, limit)
	if err != nil {
		s.logger.Error("Failed to list expense categories", "error", err)
		return nil, err
	}

	s.logger.Info("Expense categories retrieved successfully", "count", len(categories))
	return categories, nil
}

func (s *expenseCategoryService) Update(ctx context.Context, id int, req *domain.UpdateExpenseCategoryRequest) (*domain.ExpenseCategory, error) {
	s.logger.Info("Updating expense category", "id", id, "name", req.Name)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate and sanitize input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if category exists
	existingCategory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense category for update", "error", err, "id", id)
		return nil, err
	}

	// Update fields
	existingCategory.Name = name
	existingCategory.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existingCategory); err != nil {
		s.logger.Error("Failed to update expense category", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("Expense category updated successfully", "id", id, "name", name)
	return existingCategory, nil
}

func (s *expenseCategoryService) Delete(ctx context.Context, id int) error {
	s.logger.Info("Deleting expense category", "id", id)

	if id <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if category exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense category for deletion", "error", err, "id", id)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete expense category", "error", err, "id", id)
		return err
	}

	s.logger.Info("Expense category deleted successfully", "id", id)
	return nil
}
