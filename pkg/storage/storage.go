package storage

import "github.com/alex-a-renoire/sigma-homework/model"

type Storage interface {
	AddPerson(p model.Person) (int, error)
	GetPerson(id int) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id int, p model.Person) (model.Person, error)
	DeletePerson(id int) error
}
