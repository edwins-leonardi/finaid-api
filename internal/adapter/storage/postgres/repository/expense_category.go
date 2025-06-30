package repository

import (
	"context"
	"errors"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseCategoryRepository struct {
	db *pgxpool.Pool
}

// NewExpenseCategoryRepository creates a new PostgreSQL expense category repository
func NewExpenseCategoryRepository(db *pgxpool.Pool) port.ExpenseCategoryRepository {
	return &expenseCategoryRepository{
		db: db,
	}
}

func (r *expenseCategoryRepository) Create(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		INSERT INTO expense_categories (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRow(ctx, query, category.Name, category.CreatedAt, category.UpdatedAt).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *expenseCategoryRepository) GetByID(ctx context.Context, id int) (*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM expense_categories
		WHERE id = $1`

	category := &domain.ExpenseCategory{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return category, nil
}

func (r *expenseCategoryRepository) List(ctx context.Context, skip, limit int) ([]*domain.ExpenseCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM expense_categories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.ExpenseCategory
	for rows.Next() {
		category := &domain.ExpenseCategory{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *expenseCategoryRepository) Update(ctx context.Context, category *domain.ExpenseCategory) error {
	query := `
		UPDATE expense_categories
		SET name = $2, updated_at = $3
		WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, category.ID, category.Name, category.UpdatedAt)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}

func (r *expenseCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM expense_categories WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
