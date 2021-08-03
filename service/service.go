package service

import (
	"github.com/alex-a-renoire/sigma-homework/model"
)

type PersonStorage interface {
	AddPerson(name string) (int, error)
	GetPerson(id int) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id int, name string) (model.Person, error)
	DeletePerson(id int) error
}

type PersonService struct {
	db PersonStorage
}

func New(db PersonStorage) PersonService {
	return PersonService{
		db: db,
	}
}

func (s PersonService) AddPerson(name string) (int, error) {
	return s.db.AddPerson(name)
}

func (s PersonService) GetPerson(id int) (model.Person, error) {
	return s.db.GetPerson(id)
}

func (s PersonService) GetAllPersons() ([]model.Person, error) {
	return s.db.GetAllPersons()
}

func (s PersonService) UpdatePerson(id int, name string) (model.Person, error) {
	return s.db.UpdatePerson(id, name)
}

func (s PersonService) DeletePerson(id int) error {
	return s.db.DeletePerson(id)
}
