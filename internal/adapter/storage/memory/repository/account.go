package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

type AccountRepository struct {
	data map[uint64]*domain.Account
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		data: make(map[uint64]*domain.Account),
	}
}

// CreateAccount inserts a new Account into the repository
func (r *AccountRepository) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	account.ID = uint64(len(r.data) + 1) // Simple ID generation logic
	account.CreatedAt = time.Now()
	account.UpdatedAt = account.CreatedAt

	if _, exists := r.data[account.ID]; exists {
		return nil, domain.ErrConflictingData
	}
	r.data[account.ID] = account
	return account, nil
}

// GetAccountByID selects an Account by id
func (r *AccountRepository) GetAccountByID(ctx context.Context, id uint64) (*domain.Account, error) {
	account, exists := r.data[id]
	if !exists {
		return nil, domain.ErrDataNotFound
	}
	return account, nil
}

// ListAccounts selects a list of Accounts with pagination
func (r *AccountRepository) ListAccounts(ctx context.Context, skip, limit uint64) ([]domain.Account, error) {
	slog.Info("Listing accounts repo", "skip", skip, "limit", limit)
	var accounts []domain.Account
	for _, account := range r.data {
		accounts = append(accounts, *account)
	}
	slog.Info("Accounts found", "count", len(accounts))
	if skip >= uint64(len(accounts)) {
		return nil, nil // No data to return
	}
	end := skip + limit
	if end > uint64(len(accounts)) {
		end = uint64(len(accounts))
	}
	return accounts[skip:end], nil
}

// UpdateAccount updates an Account
func (r *AccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	existingAccount, exists := r.data[account.ID]
	if !exists {
		return nil, domain.ErrDataNotFound
	}
	// Update the existing account's fields
	existingAccount.Name = account.Name
	existingAccount.Currency = account.Currency
	existingAccount.AccountType = account.AccountType
	existingAccount.InitialBalance = account.InitialBalance
	existingAccount.PrimaryOwnerID = account.PrimaryOwnerID
	existingAccount.SecondOwnerID = account.SecondOwnerID
	existingAccount.UpdatedAt = time.Now()

	r.data[account.ID] = existingAccount
	return existingAccount, nil
}

// DeleteAccount deletes an Account
func (r *AccountRepository) DeleteAccount(ctx context.Context, id uint64) error {
	if _, exists := r.data[id]; !exists {
		return domain.ErrDataNotFound
	}
	delete(r.data, id)
	return nil
}
