package service

import (
	"mime/multipart"

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

type PersonService interface {
	AddPerson(name string) (uuid.UUID, error)
	GetPerson(id uuid.UUID) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id uuid.UUID, person model.Person) error
	DeletePerson(id uuid.UUID) error
	ProcessCSV(file multipart.File) error
	DownloadPersonsCSV() ([]byte, error)
}
