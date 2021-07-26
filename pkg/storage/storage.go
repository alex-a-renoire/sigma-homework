package storage

import "github.com/alex-a-renoire/tcp/model"

type Storage interface {
	AddPerson(name string) (int, error)
	GetPerson(id int) (model.Person, error)
	UpdatePerson(id int, name string) (model.Person, error)
	DeletePerson(id int) error
}
