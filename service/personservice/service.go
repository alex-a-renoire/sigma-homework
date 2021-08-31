package personservice

import (
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
)

type PersonStorage interface {
	AddPerson(p model.Person) (uuid.UUID, error)
	GetPerson(id uuid.UUID) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id uuid.UUID, person model.Person) error
	DeletePerson(id uuid.UUID) error
}

type PersonService struct {
	db PersonStorage
}

func New(db PersonStorage) PersonService {
	return PersonService{
		db: db,
	}
}

func (s PersonService) AddPerson(p model.AddUpdatePerson) (uuid.UUID, error) {
	if err := p.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to validate person: %w", err)
	}

	return s.db.AddPerson(model.Person{
		Name: p.Name,
	})
}

func (s PersonService) GetPerson(id uuid.UUID) (model.Person, error) {
	if id == uuid.Nil {
		return model.Person{}, fmt.Errorf("id should not be nil")
	}

	return s.db.GetPerson(id)
}

func (s PersonService) GetAllPersons() ([]model.Person, error) {
	return s.db.GetAllPersons()
}

func (s PersonService) UpdatePerson(id uuid.UUID, p model.AddUpdatePerson) error {
	if id == uuid.Nil {
		return fmt.Errorf("id should not be nil")
	}
	if err := p.Validate(); err != nil {
		return fmt.Errorf("failed to validate person: %w", err)
	}

	if err := s.db.UpdatePerson(id, model.Person{
		Name: p.Name,
	}); err != nil {
		return fmt.Errorf("service: failed to update person: %w", err)
	}
	return nil
}

func (s PersonService) DeletePerson(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id should not be nil")
	}

	return s.db.DeletePerson(id)
}
