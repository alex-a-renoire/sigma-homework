package service

import (
	"mime/multipart"

	"github.com/alex-a-renoire/sigma-homework/model"
)

type DierctPersonService struct {
	db PersonStorage
}

func NewDirect(db PersonStorage) DierctPersonService {
	return DierctPersonService{
		db: db,
	}
}

func (s DierctPersonService) AddPerson(name string) (int, error) {
	return s.db.AddPerson(model.Person{
		Name: name,
	})
}

func (s DierctPersonService) GetPerson(id int) (model.Person, error) {
	return s.db.GetPerson(id)
}

func (s DierctPersonService) GetAllPersons() ([]model.Person, error) {
	return s.db.GetAllPersons()
}

func (s DierctPersonService) UpdatePerson(id int, p model.Person) error {
	s.db.UpdatePerson(id, p)
	return nil
}

func (s DierctPersonService) DeletePerson(id int) error {
	return s.db.DeletePerson(id)
}

func (s DierctPersonService) DownloadPersonsCSV() ([]byte, error) {
	return nil, nil
}

func (s DierctPersonService) ProcessCSV(file multipart.File) error {
	return nil
}
