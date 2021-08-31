package storage

import (
	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
)

var _ Storage = MockStorage{}

type MockStorage struct {
	MockAddPerson     func(p model.Person) (uuid.UUID, error)
	MockGetPerson     func(id uuid.UUID) (model.Person, error)
	MockGetAllPersons func() ([]model.Person, error)
	MockUpdatePerson  func(id uuid.UUID, person model.Person) error
	MockDeletePerson  func(id uuid.UUID) error
}

func (m MockStorage) AddPerson(p model.Person) (uuid.UUID, error) {
	return m.MockAddPerson(p)
}

func (m MockStorage) GetPerson(id uuid.UUID) (model.Person, error) {
	return m.MockGetPerson(id)
}

func (m MockStorage) GetAllPersons() ([]model.Person, error) {
	return m.MockGetAllPersons()
}

func (m MockStorage) UpdatePerson(id uuid.UUID, p model.Person) error {
	return m.MockUpdatePerson(id, p)
}

func (m MockStorage) DeletePerson(id uuid.UUID) error {
	return m.MockDeletePerson(id)
}
