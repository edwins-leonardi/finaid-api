package repository

import (
	"context"
	"errors"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseSubCategoryRepository struct {
	db *pgxpool.Pool
}

// NewExpenseSubCategoryRepository creates a new PostgreSQL expense subcategory repository
func NewExpenseSubCategoryRepository(db *pgxpool.Pool) port.ExpenseSubCategoryRepository {
	return &expenseSubCategoryRepository{
		db: db,
	}
}

func (r *expenseSubCategoryRepository) Create(ctx context.Context, subcategory *domain.ExpenseSubCategory) error {
	query := `
		INSERT INTO expense_subcategories (name, expense_category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := r.db.QueryRow(ctx, query, subcategory.Name, subcategory.ExpenseCategoryID, subcategory.CreatedAt, subcategory.UpdatedAt).Scan(&subcategory.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *expenseSubCategoryRepository) GetByID(ctx context.Context, id int) (*domain.ExpenseSubCategory, error) {
	query := `
		SELECT id, name, expense_category_id, created_at, updated_at
		FROM expense_subcategories
		WHERE id = $1`

	subcategory := &domain.ExpenseSubCategory{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subcategory.ID,
		&subcategory.Name,
		&subcategory.ExpenseCategoryID,
		&subcategory.CreatedAt,
		&subcategory.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return subcategory, nil
}

func (r *expenseSubCategoryRepository) List(ctx context.Context, skip, limit int, expenseCategoryID *int) ([]*domain.ExpenseSubCategory, error) {
	var query string
	var args []interface{}

	if expenseCategoryID != nil {
		query = `
			SELECT id, name, expense_category_id, created_at, updated_at
			FROM expense_subcategories
			WHERE expense_category_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`
		args = []interface{}{*expenseCategoryID, limit, skip}
	} else {
		query = `
			SELECT id, name, expense_category_id, created_at, updated_at
			FROM expense_subcategories
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{limit, skip}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subcategories []*domain.ExpenseSubCategory
	for rows.Next() {
		subcategory := &domain.ExpenseSubCategory{}
		err := rows.Scan(
			&subcategory.ID,
			&subcategory.Name,
			&subcategory.ExpenseCategoryID,
			&subcategory.CreatedAt,
			&subcategory.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subcategories = append(subcategories, subcategory)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subcategories, nil
}

func (r *expenseSubCategoryRepository) Update(ctx context.Context, subcategory *domain.ExpenseSubCategory) error {
	query := `
		UPDATE expense_subcategories
		SET name = $2, expense_category_id = $3, updated_at = $4
		WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, subcategory.ID, subcategory.Name, subcategory.ExpenseCategoryID, subcategory.UpdatedAt)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}

func (r *expenseSubCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM expense_subcategories WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
