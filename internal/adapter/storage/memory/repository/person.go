package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
)

type PersonRepository struct {
	data map[uint64]*domain.Person
}

func NewPersonRepository() *PersonRepository {
	return &PersonRepository{
		data: make(map[uint64]*domain.Person),
	}
}

// CreatePerson inserts a new Person into the repository
func (r *PersonRepository) CreatePerson(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	person.ID = uint64(len(r.data) + 1) // Simple ID generation logic
	person.CreatedAt = time.Now()
	person.UpdatedAt = person.CreatedAt

	if _, exists := r.data[person.ID]; exists {
		return nil, domain.ErrConflictingData
	}
	r.data[person.ID] = person
	return person, nil
}

// GetPersonByID selects a Person by id
func (r *PersonRepository) GetPersonByID(ctx context.Context, id uint64) (*domain.Person, error) {
	person, exists := r.data[id]
	if !exists {
		return nil, domain.ErrDataNotFound
	}
	return person, nil
}

// GetPersonByEmail selects a Person by email
func (r *PersonRepository) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {
	for _, person := range r.data {
		if person.Email == email {
			return person, nil
		}
	}
	return nil, domain.ErrDataNotFound
}

// ListPersons selects a list of Persons with pagination
func (r *PersonRepository) ListPersons(ctx context.Context, skip, limit uint64) ([]domain.Person, error) {
	slog.Info("Listing persons repo", "skip", skip, "limit", limit)
	var persons []domain.Person
	for _, person := range r.data {
		persons = append(persons, *person)
	}
	slog.Info("Persons found", "count", len(persons))
	if skip >= uint64(len(persons)) {
		return nil, nil // No data to return
	}
	end := skip + limit
	if end > uint64(len(persons)) {
		end = uint64(len(persons))
	}
	return persons[skip:end], nil
}

// UpdatePerson updates a Person
func (r *PersonRepository) UpdatePerson(ctx context.Context, Person *domain.Person) (*domain.Person, error) {
	existingPerson, exists := r.data[Person.ID]
	if !exists {
		return nil, domain.ErrDataNotFound
	}
	// Update the existing person's fields
	existingPerson.Name = Person.Name
	existingPerson.Email = Person.Email
	// Add other fields as necessary
	r.data[Person.ID] = existingPerson
	return existingPerson, nil
}

// DeletePerson deletes a Person
func (r *PersonRepository) DeletePerson(ctx context.Context, id uint64) error {
	if _, exists := r.data[id]; !exists {
		return domain.ErrDataNotFound
	}
	delete(r.data, id)
	return nil
}
