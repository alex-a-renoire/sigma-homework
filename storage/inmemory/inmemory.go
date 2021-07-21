package inmemory

import (
	"fmt"

	dummytcp "github.com/alex-a-renoire/tcp"
)

type PersonStorage struct {
	ListPerson map[int]dummytcp.Person
	LastId     int
}

func New() PersonStorage {
	s := PersonStorage{}
	s.ListPerson = make(map[int]dummytcp.Person)
	s.LastId = 1
	return s
}

func (s *PersonStorage) AddPerson(name string) int {
	p := dummytcp.Person{
		Id:   s.LastId,
		Name: name,
	}

	s.ListPerson[s.LastId] = p
	s.LastId = s.LastId + 1

	return p.Id
}

func (s *PersonStorage) GetPerson(id int) (dummytcp.Person, error) {
	val, ok := s.ListPerson[id]
	if !ok {
		return dummytcp.Person{}, fmt.Errorf("person with id %d not found", id)
	}
	return val, nil
}

func (s *PersonStorage) UpdatePerson(id int, name string) (dummytcp.Person, error) {
	if _, ok := s.ListPerson[id]; !ok {
		return dummytcp.Person{}, fmt.Errorf("person with id %d not found", id)
	}

	s.ListPerson[id] = dummytcp.Person{
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
