package service

import (
	"context"
	"log/slog"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
)

type PersonService struct {
	repo port.PersonRepository
}

func NewPersonService(repo port.PersonRepository) *PersonService {
	return &PersonService{
		repo: repo,
	}
}

// Create created a new Person
func (svc *PersonService) Create(ctx context.Context, Person *domain.Person) (*domain.Person, error) {
	// Validate the Person data here if needed
	// For example, check if email is valid, etc.

	slog.Info("Creating new user", "name", Person.Name)

	// Call the repository to create the Person
	person, err := svc.repo.CreatePerson(ctx, Person)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return person, nil
}

// GetPerson returns a Person by id
func (svc *PersonService) GetPerson(ctx context.Context, id uint64) (*domain.Person, error) {
	person, err := svc.repo.GetPersonByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return person, nil
}

// ListPersons returns a list of Persons with pagination
func (svc *PersonService) ListPersons(ctx context.Context, skip, limit uint64) ([]domain.Person, error) {
	slog.Info("Listing persons", "skip", skip, "limit", limit)
	persons, err := svc.repo.ListPersons(ctx, skip, limit)
	slog.Info("SERVICE Persons found", "count", len(persons))
	if err != nil {
		return nil, domain.ErrInternal
	}
	return persons, nil
}

// UpdatePerson updates a Person
func (svc *PersonService) UpdatePerson(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	existingPerson, err := svc.repo.GetPersonByID(ctx, person.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	emptyData := person.Name == "" && person.Email == ""
	sameData := existingPerson.Name == person.Name && existingPerson.Email == person.Email
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	// Call the repository to update the Person
	updatedPerson, err := svc.repo.UpdatePerson(ctx, person)
	if err != nil {
		if err == domain.ErrNoUpdatedData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return updatedPerson, nil
}

// DeletePerson deletes a Person
func (svc *PersonService) DeletePerson(ctx context.Context, id uint64) error {

	_, err := svc.repo.GetPersonByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return svc.repo.DeletePerson(ctx, id)
}
