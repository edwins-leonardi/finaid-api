package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/adapter/storage/postgres"
	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	db *postgres.DB
}

func NewAccountRepository(db *postgres.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// CreateAccount inserts a new Account into the repository
func (r *AccountRepository) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	query := `
		INSERT INTO account (name, currency, account_type, initial_balance, primary_owner_id, second_owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		account.Name,
		account.Currency,
		account.AccountType,
		account.InitialBalance,
		account.PrimaryOwnerID,
		account.SecondOwnerID,
		now,
		now,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByID selects an Account by id
func (r *AccountRepository) GetAccountByID(ctx context.Context, id uint64) (*domain.Account, error) {
	query := `
		SELECT id, name, currency, account_type, initial_balance, primary_owner_id, second_owner_id, created_at, updated_at
		FROM account
		WHERE id = $1
	`

	account := &domain.Account{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Name,
		&account.Currency,
		&account.AccountType,
		&account.InitialBalance,
		&account.PrimaryOwnerID,
		&account.SecondOwnerID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return account, nil
}

// ListAccounts selects a list of Accounts with pagination
func (r *AccountRepository) ListAccounts(ctx context.Context, skip, limit uint64) ([]domain.Account, error) {
	slog.Info("Listing accounts repo", "skip", skip, "limit", limit)

	query := `
		SELECT id, name, currency, account_type, initial_balance, primary_owner_id, second_owner_id, created_at, updated_at
		FROM account
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Currency,
			&account.AccountType,
			&account.InitialBalance,
			&account.PrimaryOwnerID,
			&account.SecondOwnerID,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	slog.Info("Accounts found", "count", len(accounts))
	return accounts, nil
}

// UpdateAccount updates an Account
func (r *AccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	query := `
		UPDATE account
		SET name = $1, currency = $2, account_type = $3, initial_balance = $4, primary_owner_id = $5, second_owner_id = $6, updated_at = $7
		WHERE id = $8
		RETURNING id, name, currency, account_type, initial_balance, primary_owner_id, second_owner_id, created_at, updated_at
	`

	now := time.Now()
	updatedAccount := &domain.Account{}
	err := r.db.QueryRow(ctx, query,
		account.Name,
		account.Currency,
		account.AccountType,
		account.InitialBalance,
		account.PrimaryOwnerID,
		account.SecondOwnerID,
		now,
		account.ID,
	).Scan(
		&updatedAccount.ID,
		&updatedAccount.Name,
		&updatedAccount.Currency,
		&updatedAccount.AccountType,
		&updatedAccount.InitialBalance,
		&updatedAccount.PrimaryOwnerID,
		&updatedAccount.SecondOwnerID,
		&updatedAccount.CreatedAt,
		&updatedAccount.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return updatedAccount, nil
}

// DeleteAccount deletes an Account
func (r *AccountRepository) DeleteAccount(ctx context.Context, id uint64) error {
	query := `DELETE FROM account WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
