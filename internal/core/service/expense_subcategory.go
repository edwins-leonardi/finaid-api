package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseSubCategoryService struct {
	repo         port.ExpenseSubCategoryRepository
	categoryRepo port.ExpenseCategoryRepository
	logger       *slog.Logger
}

// NewExpenseSubCategoryService creates a new expense subcategory service
func NewExpenseSubCategoryService(
	repo port.ExpenseSubCategoryRepository,
	categoryRepo port.ExpenseCategoryRepository,
	logger *slog.Logger,
) port.ExpenseSubCategoryService {
	return &expenseSubCategoryService{
		repo:         repo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *expenseSubCategoryService) Create(ctx context.Context, req *domain.CreateExpenseSubCategoryRequest) (*domain.ExpenseSubCategory, error) {
	s.logger.Info("Creating expense subcategory", "name", req.Name, "expense_category_id", req.ExpenseCategoryID)

	// Validate and sanitize input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	if req.ExpenseCategoryID <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate that the expense category exists
	_, err := s.categoryRepo.GetByID(ctx, req.ExpenseCategoryID)
	if err != nil {
		s.logger.Error("Expense category not found", "error", err, "expense_category_id", req.ExpenseCategoryID)
		return nil, err
	}

	subcategory := &domain.ExpenseSubCategory{
		Name:              name,
		ExpenseCategoryID: req.ExpenseCategoryID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(ctx, subcategory); err != nil {
		s.logger.Error("Failed to create expense subcategory", "error", err, "name", name)
		return nil, err
	}

	s.logger.Info("Expense subcategory created successfully", "id", subcategory.ID, "name", subcategory.Name)
	return subcategory, nil
}

func (s *expenseSubCategoryService) GetByID(ctx context.Context, id int) (*domain.ExpenseSubCategory, error) {
	s.logger.Info("Getting expense subcategory by ID", "id", id)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	subcategory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense subcategory", "error", err, "id", id)
		return nil, err
	}

	return subcategory, nil
}

func (s *expenseSubCategoryService) List(ctx context.Context, req *domain.ListExpenseSubCategoriesRequest) ([]*domain.ExpenseSubCategory, error) {
	s.logger.Info("Listing expense subcategories", "skip", req.Skip, "limit", req.Limit, "expense_category_id", req.ExpenseCategoryID)

	// Set default values
	skip := req.Skip
	if skip < 0 {
		skip = 0
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Optional filter by expense category
	var expenseCategoryID *int
	if req.ExpenseCategoryID > 0 {
		expenseCategoryID = &req.ExpenseCategoryID
	}

	subcategories, err := s.repo.List(ctx, skip, limit, expenseCategoryID)
	if err != nil {
		s.logger.Error("Failed to list expense subcategories", "error", err)
		return nil, err
	}

	s.logger.Info("Expense subcategories retrieved successfully", "count", len(subcategories))
	return subcategories, nil
}

func (s *expenseSubCategoryService) Update(ctx context.Context, id int, req *domain.UpdateExpenseSubCategoryRequest) (*domain.ExpenseSubCategory, error) {
	s.logger.Info("Updating expense subcategory", "id", id, "name", req.Name, "expense_category_id", req.ExpenseCategoryID)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate and sanitize input
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	if req.ExpenseCategoryID <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Check if subcategory exists
	existingSubCategory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense subcategory for update", "error", err, "id", id)
		return nil, err
	}

	// Validate that the expense category exists
	_, err = s.categoryRepo.GetByID(ctx, req.ExpenseCategoryID)
	if err != nil {
		s.logger.Error("Expense category not found", "error", err, "expense_category_id", req.ExpenseCategoryID)
		return nil, err
	}

	// Update fields
	existingSubCategory.Name = name
	existingSubCategory.ExpenseCategoryID = req.ExpenseCategoryID
	existingSubCategory.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existingSubCategory); err != nil {
		s.logger.Error("Failed to update expense subcategory", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("Expense subcategory updated successfully", "id", id, "name", name)
	return existingSubCategory, nil
}

func (s *expenseSubCategoryService) Delete(ctx context.Context, id int) error {
	s.logger.Info("Deleting expense subcategory", "id", id)

	if id <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if subcategory exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense subcategory for deletion", "error", err, "id", id)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete expense subcategory", "error", err, "id", id)
		return err
	}

	s.logger.Info("Expense subcategory deleted successfully", "id", id)
	return nil
}
