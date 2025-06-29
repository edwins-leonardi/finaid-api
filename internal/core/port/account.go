package port

import (
	"context"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

type AccountRepository interface {
	// CreateAccount inserts a new Account into the database
	CreateAccount(ctx context.Context, Account *domain.Account) (*domain.Account, error)
	// GetAccountByID selects a Account by id
	GetAccountByID(ctx context.Context, id uint64) (*domain.Account, error)
	// ListAccounts selects a list of Accounts with pagination
	ListAccounts(ctx context.Context, skip, limit uint64) ([]domain.Account, error)
	// UpdateAccount updates a Account
	UpdateAccount(ctx context.Context, Account *domain.Account) (*domain.Account, error)
	// DeleteAccount deletes a Account
	DeleteAccount(ctx context.Context, id uint64) error
}

// AccountService is an interface for interacting with Account-related business logic
type AccountService interface {
	// Create creates a new Account
	Create(ctx context.Context, Account *domain.Account) (*domain.Account, error)
	// GetAccount returns a Account by id
	GetAccount(ctx context.Context, id uint64) (*domain.Account, error)
	// ListAccounts returns a list of Accounts with pagination
	ListAccounts(ctx context.Context, skip, limit uint64) ([]domain.Account, error)
	// UpdateAccount updates a Account
	UpdateAccount(ctx context.Context, Account *domain.Account) (*domain.Account, error)
	// DeleteAccount deletes a Account
	DeleteAccount(ctx context.Context, id uint64) error
}
