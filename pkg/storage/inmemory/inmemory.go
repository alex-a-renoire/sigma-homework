package inmemory

import (
	"fmt"

	"github.com/alex-a-renoire/tcp/model"
)

type PersonStorage struct {
	ListPerson map[int]model.Person
	LastId     int
}

func New() PersonStorage {
	s := PersonStorage{}
	s.ListPerson = make(map[int]model.Person)
	s.LastId = 1
	return s
}

func (s *PersonStorage) AddPerson(name string) int {
	p := model.Person{
		Id:   s.LastId,
		Name: name,
	}

	s.ListPerson[s.LastId] = p
	s.LastId = s.LastId + 1

	return p.Id
}

func (s *PersonStorage) GetPerson(id int) (model.Person, error) {
	val, ok := s.ListPerson[id]
	if !ok {
		return model.Person{}, fmt.Errorf("person with id %d not found", id)
	}
	return val, nil
}

func (s *PersonStorage) UpdatePerson(id int, name string) (model.Person, error) {
	if _, ok := s.ListPerson[id]; !ok {
		return model.Person{}, fmt.Errorf("person with id %d not found", id)
	}

	s.ListPerson[id] = model.Person{
		Id:   id,
		Name: name,
	}
	return s.ListPerson[id], nil
}

func (s *PersonStorage) DeletePerson(id int) error {
	if _, ok := s.ListPerson[id]; !ok {
		return fmt.Errorf("person with id %d not found", id)
	}

	delete(s.ListPerson, id)
	return nil
}
