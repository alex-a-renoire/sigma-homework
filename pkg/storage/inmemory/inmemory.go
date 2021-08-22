package inmemory

import (
	"fmt"
	"sync"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
)

type PersonStorage struct {
	MapPerson map[uuid.UUID]model.Person
	mu        sync.Mutex
}

func New() *PersonStorage {
	s := PersonStorage{}
	s.MapPerson = make(map[uuid.UUID]model.Person)
	return &s
}

func (s *PersonStorage) AddPerson(p model.Person) (uuid.UUID, error) {
	p.Id = uuid.New()
	s.mu.Lock()
	s.MapPerson[p.Id] = p
	s.mu.Unlock()

	return p.Id, nil
}

func (s *PersonStorage) GetPerson(id uuid.UUID) (model.Person, error) {
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

func (s *PersonStorage) UpdatePerson(id uuid.UUID, p model.Person) error {
	s.mu.Lock()
	if _, ok := s.MapPerson[id]; !ok {
		return fmt.Errorf("person with id %d not found", id)
	}

	p.Id = id
	s.MapPerson[id] = p
	s.mu.Unlock()
	return nil
}

func (s *PersonStorage) DeletePerson(id uuid.UUID) error {
	s.mu.Lock()
	if _, ok := s.MapPerson[id]; !ok {
		return fmt.Errorf("person with id %d not found", id)
	}

	delete(s.MapPerson, id)
	s.mu.Unlock()
	return nil
}
