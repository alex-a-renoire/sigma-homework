package service

import (
	"github.com/alex-a-renoire/sigma-homework/model"
)

type DierctPersonService struct {
	db PersonStorage
}

func New(db PersonStorage) DierctPersonService {
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

func (s DierctPersonService) UpdatePerson(id int, name string) (model.Person, error) {
	return s.db.UpdatePerson(id, model.Person{
		Name: name,
	})
}

func (s DierctPersonService) DeletePerson(id int) error {
	return s.db.DeletePerson(id)
}
