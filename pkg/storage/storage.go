package storage

import (
	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
)

//TODO save by UUID

type Storage interface {
	AddPerson(p model.Person) (uuid.UUID, error)
	GetPerson(id uuid.UUID) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id uuid.UUID, p model.Person) error
	DeletePerson(id uuid.UUID) error
}
