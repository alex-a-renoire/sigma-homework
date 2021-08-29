package storage

import (
	"mime/multipart"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
)

type Storage interface {
	AddPerson(p model.Person) (uuid.UUID, error)
	GetPerson(id uuid.UUID) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id uuid.UUID, p model.Person) error
	DeletePerson(id uuid.UUID) error
}

type CsvProcessor interface {
	DownloadPersonsCSV() ([]byte, error)
	ProcessCSV(file multipart.File) error
}
