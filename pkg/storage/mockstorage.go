package storage

import "github.com/alex-a-renoire/tcp/model"

var _ Storage = MockStorage{}

type MockStorage struct {
	MockAddPerson    func(name string) (int, error)
	MockGetPerson    func(id int) (model.Person, error)
	MockUpdatePerson func(id int, name string) (model.Person, error)
	MockDeletePerson func(id int) error
}

func (m MockStorage) AddPerson(name string) (int, error) {
	return m.MockAddPerson(name)
}

func (m MockStorage) GetPerson(id int) (model.Person, error) {
	return m.MockGetPerson(id)
}

func (m MockStorage) UpdatePerson(id int, name string) (model.Person, error) {
	return m.MockUpdatePerson(id, name)
}

func (m MockStorage) DeletePerson(id int) error {
	return m.MockDeletePerson(id)
}