package service

import "github.com/alex-a-renoire/sigma-homework/model"

type PersonStorage interface {
	AddPerson(name string) (int, error)
	GetPerson(id int) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id int, name string) (model.Person, error)
	DeletePerson(id int) error
}

type PersonService interface {
	AddPerson(name string) (int, error) 
	GetPerson(id int) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id int, name string) (model.Person, error) 
	DeletePerson(id int) error
}