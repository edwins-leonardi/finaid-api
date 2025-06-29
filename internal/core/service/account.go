package service

import (
	"context"
	"log/slog"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type AccountService struct {
	repo       port.AccountRepository
	personRepo port.PersonRepository
}

func NewAccountService(repo port.AccountRepository, personRepo port.PersonRepository) *AccountService {
	return &AccountService{
		repo:       repo,
		personRepo: personRepo,
	}
}

// Create creates a new Account
func (svc *AccountService) Create(ctx context.Context, Account *domain.Account) (*domain.Account, error) {
	// Validate the Account data here if needed
	// For example, check if account type is valid, currency format, etc.

	// Validate that primary owner exists
	_, err := svc.personRepo.GetPersonByID(ctx, Account.PrimaryOwnerID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			slog.Error("Primary owner not found", "primary_owner_id", Account.PrimaryOwnerID)
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	// Validate that second owner exists (if provided)
	if Account.SecondOwnerID != nil {
		_, err := svc.personRepo.GetPersonByID(ctx, *Account.SecondOwnerID)
		if err != nil {
			if err == domain.ErrDataNotFound {
				slog.Error("Second owner not found", "second_owner_id", *Account.SecondOwnerID)
				return nil, domain.ErrDataNotFound
			}
			return nil, domain.ErrInternal
		}
	}

	slog.Info("Creating new account", "name", Account.Name, "type", Account.AccountType, "primary_owner_id", Account.PrimaryOwnerID)

	// Call the repository to create the Account
	account, err := svc.repo.CreateAccount(ctx, Account)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return account, nil
}

// GetAccount returns a Account by id
func (svc *AccountService) GetAccount(ctx context.Context, id uint64) (*domain.Account, error) {
	account, err := svc.repo.GetAccountByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return account, nil
}

// ListAccounts returns a list of Accounts with pagination
func (svc *AccountService) ListAccounts(ctx context.Context, skip, limit uint64) ([]domain.Account, error) {
	slog.Info("Listing accounts", "skip", skip, "limit", limit)
	accounts, err := svc.repo.ListAccounts(ctx, skip, limit)
	slog.Info("SERVICE Accounts found", "count", len(accounts))
	if err != nil {
		return nil, domain.ErrInternal
	}
	return accounts, nil
}

// UpdateAccount updates a Account
func (svc *AccountService) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	existingAccount, err := svc.repo.GetAccountByID(ctx, account.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// Validate that primary owner exists
	_, err = svc.personRepo.GetPersonByID(ctx, account.PrimaryOwnerID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			slog.Error("Primary owner not found", "primary_owner_id", account.PrimaryOwnerID)
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	// Validate that second owner exists (if provided)
	if account.SecondOwnerID != nil {
		_, err := svc.personRepo.GetPersonByID(ctx, *account.SecondOwnerID)
		if err != nil {
			if err == domain.ErrDataNotFound {
				slog.Error("Second owner not found", "second_owner_id", *account.SecondOwnerID)
				return nil, domain.ErrDataNotFound
			}
			return nil, domain.ErrInternal
		}
	}

	// Check if there's actually something to update
	emptyData := account.Name == "" && account.Currency == "" && account.AccountType == "" && account.InitialBalance == 0
	sameData := existingAccount.Name == account.Name &&
		existingAccount.Currency == account.Currency &&
		existingAccount.AccountType == account.AccountType &&
		existingAccount.InitialBalance == account.InitialBalance &&
		existingAccount.PrimaryOwnerID == account.PrimaryOwnerID &&
		((existingAccount.SecondOwnerID == nil && account.SecondOwnerID == nil) ||
			(existingAccount.SecondOwnerID != nil && account.SecondOwnerID != nil && *existingAccount.SecondOwnerID == *account.SecondOwnerID))

	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	// Call the repository to update the Account
	updatedAccount, err := svc.repo.UpdateAccount(ctx, account)
	if err != nil {
		if err == domain.ErrNoUpdatedData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedAccount, nil
}

// DeleteAccount deletes a Account
func (svc *AccountService) DeleteAccount(ctx context.Context, id uint64) error {
	_, err := svc.repo.GetAccountByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return svc.repo.DeleteAccount(ctx, id)
}
