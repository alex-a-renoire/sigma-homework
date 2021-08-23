package redisstorage

import (
	"encoding/json"
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RDSdb struct {
	Client *redis.Client
}

func NewRDS(addr string, pwd string, dbname int) *RDSdb {
	return &RDSdb{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       dbname,
		}),
	}
}

func (db *RDSdb) AddPerson(p model.Person) (uuid.UUID, error) {
	p.Id = uuid.New()

	person, err := json.Marshal(p)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Cannot add person to db: %w", err)
	}

	_, err = db.Client.Set("person:"+p.Id.String(), person, 0).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Cannot add person to db: %w", err)
	}

	return p.Id, nil
}

func (db *RDSdb) GetPerson(id uuid.UUID) (model.Person, error) {
	p, err := db.Client.Get("person:" + id.String()).Result()
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to find person: %w", err)
	}

	var person model.Person

	if err := json.Unmarshal([]byte(p), &person); err != nil {
		return model.Person{}, fmt.Errorf("Failed to marshal persons string: %w", err)
	}

	return person, nil
}

func (db *RDSdb) GetAllPersons() ([]model.Person, error) {
	var persons []model.Person

	keys, err := db.Client.Keys("person:*").Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve persons from db: %w", err)
	}

	for _, k := range keys {
		var person model.Person
		reply, err := db.Client.Get(k).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve person by key from db: %w", err)
		}

		if err := json.Unmarshal([]byte(reply), &person); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal persons string: %w", err)
		}

		// person.Id, err = uuid.Parse(strings.TrimPrefix(k, "person:"))
		// if err != nil {
		// 	return nil, fmt.Errorf("malformed id or prefix: %w", err)
		// }

		persons = append(persons, person)
	}

	return persons, nil
}

func (db *RDSdb) UpdatePerson(id uuid.UUID, p model.Person) error {
	person, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Cannot marshal person: %w", err)
	}

	_, err = db.Client.Set("person:"+id.String(), person, 0).Result()
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}

func (db *RDSdb) DeletePerson(id uuid.UUID) error {
	_, err := db.Client.Del("person:" + id.String()).Result()
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}
