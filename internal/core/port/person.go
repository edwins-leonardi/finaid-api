package port

import (
	"context"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

type PersonRepository interface {
	// CreatePerson inserts a new Person into the database
	CreatePerson(ctx context.Context, Person *domain.Person) (*domain.Person, error)
	// GetPersonByID selects a Person by id
	GetPersonByID(ctx context.Context, id uint64) (*domain.Person, error)
	// GetPersonByEmail selects a Person by email
	GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error)
	// ListPersons selects a list of Persons with pagination
	ListPersons(ctx context.Context, skip, limit uint64) ([]domain.Person, error)
	// UpdatePerson updates a Person
	UpdatePerson(ctx context.Context, Person *domain.Person) (*domain.Person, error)
	// DeletePerson deletes a Person
	DeletePerson(ctx context.Context, id uint64) error
}

// PersonService is an interface for interacting with Person-related business logic
type PersonService interface {
	// Create creates a new Person
	Create(ctx context.Context, Person *domain.Person) (*domain.Person, error)
	// GetPerson returns a Person by id
	GetPerson(ctx context.Context, id uint64) (*domain.Person, error)
	// ListPersons returns a list of Persons with pagination
	ListPersons(ctx context.Context, skip, limit uint64) ([]domain.Person, error)
	// UpdatePerson updates a Person
	UpdatePerson(ctx context.Context, Person *domain.Person) (*domain.Person, error)
	// DeletePerson deletes a Person
	DeletePerson(ctx context.Context, id uint64) error
}
