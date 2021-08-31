package personservice

import (
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

func (s PersonService) AddPerson(p model.Person) (uuid.UUID, error) {
	p.Validate()
	return s.db.AddPerson(model.Person{
		Name: p.Name,
	})
}

func (s PersonService) GetPerson(id uuid.UUID) (model.Person, error) {
	return s.db.GetPerson(id)
}

func (s PersonService) GetAllPersons() ([]model.Person, error) {
	return s.db.GetAllPersons()
}

func (s PersonService) UpdatePerson(id uuid.UUID, p model.Person) error {
	s.db.UpdatePerson(id, p)
	return nil
}

func (s PersonService) DeletePerson(id uuid.UUID) error {
	return s.db.DeletePerson(id)
}
