package inmemory

import (
	dummytcp "github.com/alex-a-renoire/tcp"
)

type PersonStorage struct {
	ListPerson map[int]dummytcp.Person
	LastId     int
}

func New() PersonStorage {
	s := PersonStorage{}
	s.ListPerson = make(map[int]dummytcp.Person)
	s.LastId = 0
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

func (s *PersonStorage) GetPerson(id int) dummytcp.Person {
	return s.ListPerson[id]
}

func (s *PersonStorage) UpdatePerson(id int, name string) dummytcp.Person {
	s.ListPerson[id] = dummytcp.Person{
		Id:   id,
		Name: name,
	}
	return s.ListPerson[id]
}

func (s *PersonStorage) DeletePerson(id int) {
	_, ok := s.ListPerson[id]
	if ok {
		delete(s.ListPerson, id)
	}
}
