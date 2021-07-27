package inmemory

import (
	"fmt"
	"sync"

	"github.com/alex-a-renoire/sigma-homework/model"
)

type PersonStorage struct {
	MapPerson map[int]model.Person
	LastId    int
	mu        sync.Mutex
}

func New() *PersonStorage {
	s := PersonStorage{}
	s.MapPerson = make(map[int]model.Person)
	s.LastId = 1
	return &s
}

func (s *PersonStorage) AddPerson(name string) (int, error) {
	p := model.Person{
		Id:   s.LastId,
		Name: name,
	}

	s.mu.Lock()
	s.MapPerson[s.LastId] = p
	s.LastId = s.LastId + 1
	s.mu.Unlock()

	return p.Id, nil
}

func (s *PersonStorage) GetPerson(id int) (model.Person, error) {
	s.mu.Lock()
	val, ok := s.MapPerson[id]
	s.mu.Unlock()

	if !ok {
		return model.Person{}, fmt.Errorf("person with id %d not found", id)
	}
	return val, nil
}

func (s *PersonStorage) GetAllPersons() ([]model.Person, error) {
	s.mu.Lock()
	persons := make([]model.Person, 0, len(s.MapPerson))

	for _, person := range s.MapPerson {
		persons = append(persons, person)
	}
	s.mu.Unlock()
	return persons, nil
}

func (s *PersonStorage) UpdatePerson(id int, name string) (model.Person, error) {
	s.mu.Lock()
	if _, ok := s.MapPerson[id]; !ok {
		return model.Person{}, fmt.Errorf("person with id %d not found", id)
	}

	s.MapPerson[id] = model.Person{
		Id:   id,
		Name: name,
	}
	s.mu.Unlock()
	return s.MapPerson[id], nil
}

func (s *PersonStorage) DeletePerson(id int) error {
	s.mu.Lock()
	if _, ok := s.MapPerson[id]; !ok {
		return fmt.Errorf("person with id %d not found", id)
	}

	delete(s.MapPerson, id)
	s.mu.Unlock()
	return nil
}
