package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/adapter/storage/postgres"
	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type PersonRepository struct {
	db *postgres.DB
}

func NewPersonRepository(db *postgres.DB) *PersonRepository {
	return &PersonRepository{
		db: db,
	}
}

// CreatePerson inserts a new Person into the repository
func (r *PersonRepository) CreatePerson(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	query := `
		INSERT INTO person (name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		person.Name,
		person.Email,
		now,
		now,
	).Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return person, nil
}

// GetPersonByID selects a Person by id
func (r *PersonRepository) GetPersonByID(ctx context.Context, id uint64) (*domain.Person, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM person
		WHERE id = $1
	`

	person := &domain.Person{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Email,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return person, nil
}

// GetPersonByEmail selects a Person by email
func (r *PersonRepository) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM person
		WHERE email = $1
	`

	person := &domain.Person{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&person.ID,
		&person.Name,
		&person.Email,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return person, nil
}

// ListPersons selects a list of Persons with pagination
func (r *PersonRepository) ListPersons(ctx context.Context, skip, limit uint64) ([]domain.Person, error) {
	slog.Info("Listing persons repo", "skip", skip, "limit", limit)

	query := `
		SELECT id, name, email, created_at, updated_at
		FROM person
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []domain.Person
	for rows.Next() {
		var person domain.Person
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Email,
			&person.CreatedAt,
			&person.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	slog.Info("Persons found", "count", len(persons))
	return persons, nil
}

// UpdatePerson updates a Person
func (r *PersonRepository) UpdatePerson(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	query := `
		UPDATE person
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, email, created_at, updated_at
	`

	now := time.Now()
	updatedPerson := &domain.Person{}
	err := r.db.QueryRow(ctx, query,
		person.Name,
		person.Email,
		now,
		person.ID,
	).Scan(
		&updatedPerson.ID,
		&updatedPerson.Name,
		&updatedPerson.Email,
		&updatedPerson.CreatedAt,
		&updatedPerson.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return updatedPerson, nil
}

// DeletePerson deletes a Person
func (r *PersonRepository) DeletePerson(ctx context.Context, id uint64) error {
	query := `DELETE FROM person WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
