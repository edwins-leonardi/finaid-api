package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseRepository struct {
	db *pgxpool.Pool
}

// NewExpenseRepository creates a new PostgreSQL expense repository
func NewExpenseRepository(db *pgxpool.Pool) port.ExpenseRepository {
	return &expenseRepository{
		db: db,
	}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (amount, category_id, subcategory_id, date, payee_id, account_id, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.db.QueryRow(ctx, query,
		expense.Amount,
		expense.CategoryID,
		expense.SubCategoryID,
		expense.Date,
		expense.PayeeID,
		expense.AccountID,
		expense.Notes,
		expense.CreatedAt,
		expense.UpdatedAt,
	).Scan(&expense.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *expenseRepository) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	query := `
		SELECT id, amount, category_id, subcategory_id, date, payee_id, account_id, notes, created_at, updated_at
		FROM expenses
		WHERE id = $1`

	expense := &domain.Expense{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&expense.ID,
		&expense.Amount,
		&expense.CategoryID,
		&expense.SubCategoryID,
		&expense.Date,
		&expense.PayeeID,
		&expense.AccountID,
		&expense.Notes,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return expense, nil
}

func (r *expenseRepository) List(ctx context.Context, filters port.ExpenseFilters) ([]*domain.Expense, error) {
	// Build dynamic query based on filters
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT id, amount, category_id, subcategory_id, date, payee_id, account_id, notes, created_at, updated_at
		FROM expenses`

	// Add WHERE conditions based on filters
	if filters.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.SubCategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("subcategory_id = $%d", argIndex))
		args = append(args, *filters.SubCategoryID)
		argIndex++
	}

	if filters.PayeeID != nil {
		conditions = append(conditions, fmt.Sprintf("payee_id = $%d", argIndex))
		args = append(args, *filters.PayeeID)
		argIndex++
	}

	if filters.AccountID != nil {
		conditions = append(conditions, fmt.Sprintf("account_id = $%d", argIndex))
		args = append(args, *filters.AccountID)
		argIndex++
	}

	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	// Build final query
	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY date DESC, created_at DESC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Skip)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.Amount,
			&expense.CategoryID,
			&expense.SubCategoryID,
			&expense.Date,
			&expense.PayeeID,
			&expense.AccountID,
			&expense.Notes,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses
		SET amount = $2, category_id = $3, subcategory_id = $4, date = $5, payee_id = $6, account_id = $7, notes = $8, updated_at = $9
		WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query,
		expense.ID,
		expense.Amount,
		expense.CategoryID,
		expense.SubCategoryID,
		expense.Date,
		expense.PayeeID,
		expense.AccountID,
		expense.Notes,
		expense.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}

func (r *expenseRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM expenses WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
