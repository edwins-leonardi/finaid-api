package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type expenseService struct {
	repo            port.ExpenseRepository
	categoryRepo    port.ExpenseCategoryRepository
	subCategoryRepo port.ExpenseSubCategoryRepository
	personRepo      port.PersonRepository
	accountRepo     port.AccountRepository
	logger          *slog.Logger
}

// NewExpenseService creates a new expense service
func NewExpenseService(
	repo port.ExpenseRepository,
	categoryRepo port.ExpenseCategoryRepository,
	subCategoryRepo port.ExpenseSubCategoryRepository,
	personRepo port.PersonRepository,
	accountRepo port.AccountRepository,
	logger *slog.Logger,
) port.ExpenseService {
	return &expenseService{
		repo:            repo,
		categoryRepo:    categoryRepo,
		subCategoryRepo: subCategoryRepo,
		personRepo:      personRepo,
		accountRepo:     accountRepo,
		logger:          logger,
	}
}

func (s *expenseService) Create(ctx context.Context, req *domain.CreateExpenseRequest) (*domain.Expense, error) {
	s.logger.Info("Creating expense", "amount", req.Amount, "category_id", req.CategoryID, "payee_id", req.PayeeID)

	// Validate amount
	if req.Amount < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate and parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		s.logger.Error("Invalid date format", "error", err, "date", req.Date)
		return nil, domain.ErrInvalidInput
	}

	// Validate that the expense category exists
	_, err = s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		s.logger.Error("Expense category not found", "error", err, "category_id", req.CategoryID)
		return nil, err
	}

	// Validate subcategory if provided
	if req.SubCategoryID != nil {
		subCategory, err := s.subCategoryRepo.GetByID(ctx, *req.SubCategoryID)
		if err != nil {
			s.logger.Error("Expense subcategory not found", "error", err, "subcategory_id", *req.SubCategoryID)
			return nil, err
		}

		// Ensure subcategory belongs to the specified category
		if subCategory.ExpenseCategoryID != req.CategoryID {
			s.logger.Error("Subcategory does not belong to the specified category",
				"subcategory_id", *req.SubCategoryID, "category_id", req.CategoryID,
				"subcategory_category_id", subCategory.ExpenseCategoryID)
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate that the payee (person) exists
	_, err = s.personRepo.GetPersonByID(ctx, uint64(req.PayeeID))
	if err != nil {
		s.logger.Error("Payee not found", "error", err, "payee_id", req.PayeeID)
		return nil, err
	}

	// Validate that the account exists
	_, err = s.accountRepo.GetAccountByID(ctx, uint64(req.AccountID))
	if err != nil {
		s.logger.Error("Account not found", "error", err, "account_id", req.AccountID)
		return nil, err
	}

	// Sanitize notes
	notes := strings.TrimSpace(req.Notes)

	expense := &domain.Expense{
		Amount:        req.Amount,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		Date:          date,
		PayeeID:       req.PayeeID,
		AccountID:     req.AccountID,
		Notes:         notes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, expense); err != nil {
		s.logger.Error("Failed to create expense", "error", err)
		return nil, err
	}

	s.logger.Info("Expense created successfully", "id", expense.ID, "amount", expense.Amount)
	return expense, nil
}

func (s *expenseService) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	s.logger.Info("Getting expense by ID", "id", id)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	expense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense", "error", err, "id", id)
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) List(ctx context.Context, req *domain.ListExpensesRequest) ([]*domain.Expense, error) {
	s.logger.Info("Listing expenses", "skip", req.Skip, "limit", req.Limit)

	// Set default values
	skip := req.Skip
	if skip < 0 {
		skip = 0
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Build filters
	filters := port.ExpenseFilters{
		Skip:  skip,
		Limit: limit,
	}

	// Optional filters
	if req.CategoryID > 0 {
		filters.CategoryID = &req.CategoryID
	}
	if req.SubCategoryID > 0 {
		filters.SubCategoryID = &req.SubCategoryID
	}
	if req.PayeeID > 0 {
		filters.PayeeID = &req.PayeeID
	}
	if req.AccountID > 0 {
		filters.AccountID = &req.AccountID
	}

	// Parse date filters
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			s.logger.Error("Invalid start date format", "error", err, "start_date", req.StartDate)
			return nil, domain.ErrInvalidInput
		}
		filters.StartDate = &startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			s.logger.Error("Invalid end date format", "error", err, "end_date", req.EndDate)
			return nil, domain.ErrInvalidInput
		}
		// Set end date to end of day
		endOfDay := endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		filters.EndDate = &endOfDay
	}

	expenses, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list expenses", "error", err)
		return nil, err
	}

	s.logger.Info("Expenses retrieved successfully", "count", len(expenses))
	return expenses, nil
}

func (s *expenseService) Update(ctx context.Context, id int, req *domain.UpdateExpenseRequest) (*domain.Expense, error) {
	s.logger.Info("Updating expense", "id", id, "amount", req.Amount, "category_id", req.CategoryID)

	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate amount
	if req.Amount < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Validate and parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		s.logger.Error("Invalid date format", "error", err, "date", req.Date)
		return nil, domain.ErrInvalidInput
	}

	// Check if expense exists
	existingExpense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense for update", "error", err, "id", id)
		return nil, err
	}

	// Validate that the expense category exists
	_, err = s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		s.logger.Error("Expense category not found", "error", err, "category_id", req.CategoryID)
		return nil, err
	}

	// Validate subcategory if provided
	if req.SubCategoryID != nil {
		subCategory, err := s.subCategoryRepo.GetByID(ctx, *req.SubCategoryID)
		if err != nil {
			s.logger.Error("Expense subcategory not found", "error", err, "subcategory_id", *req.SubCategoryID)
			return nil, err
		}

		// Ensure subcategory belongs to the specified category
		if subCategory.ExpenseCategoryID != req.CategoryID {
			s.logger.Error("Subcategory does not belong to the specified category",
				"subcategory_id", *req.SubCategoryID, "category_id", req.CategoryID,
				"subcategory_category_id", subCategory.ExpenseCategoryID)
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate that the payee (person) exists
	_, err = s.personRepo.GetPersonByID(ctx, uint64(req.PayeeID))
	if err != nil {
		s.logger.Error("Payee not found", "error", err, "payee_id", req.PayeeID)
		return nil, err
	}

	// Validate that the account exists
	_, err = s.accountRepo.GetAccountByID(ctx, uint64(req.AccountID))
	if err != nil {
		s.logger.Error("Account not found", "error", err, "account_id", req.AccountID)
		return nil, err
	}

	// Sanitize notes
	notes := strings.TrimSpace(req.Notes)

	// Update fields
	existingExpense.Amount = req.Amount
	existingExpense.CategoryID = req.CategoryID
	existingExpense.SubCategoryID = req.SubCategoryID
	existingExpense.Date = date
	existingExpense.PayeeID = req.PayeeID
	existingExpense.AccountID = req.AccountID
	existingExpense.Notes = notes
	existingExpense.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existingExpense); err != nil {
		s.logger.Error("Failed to update expense", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("Expense updated successfully", "id", id, "amount", req.Amount)
	return existingExpense, nil
}

func (s *expenseService) Delete(ctx context.Context, id int) error {
	s.logger.Info("Deleting expense", "id", id)

	if id <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if expense exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get expense for deletion", "error", err, "id", id)
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete expense", "error", err, "id", id)
		return err
	}

	s.logger.Info("Expense deleted successfully", "id", id)
	return nil
}
