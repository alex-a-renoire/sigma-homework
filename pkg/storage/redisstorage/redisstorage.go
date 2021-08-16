package redisstorage

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/go-redis/redis"
)

type RDSdb struct {
	currentPersonId int
	Client          *redis.Client
}

func NewRDS(addr string, pwd string, dbname int) *RDSdb {
	return &RDSdb{
		currentPersonId: 0,
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       dbname,
		}),
	}
}

func (db *RDSdb) AddPerson(p model.Person) (int, error) {
	db.currentPersonId += 1
	p.Id = db.currentPersonId

	person, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("Cannot add person to db: %w", err)
	}

	_, err = db.Client.Do("SET", "person:"+strconv.Itoa(p.Id), person).Result()
	if err != nil {
		return 0, fmt.Errorf("Cannot add person to db: %w", err)
	}

	return p.Id, nil
}

func (db *RDSdb) GetPerson(id int) (model.Person, error) {
	var person model.Person

	res, err := db.Client.Do("GET", "person:"+strconv.Itoa(id)).Result()

	if err = json.Unmarshal(res.([]byte), &person); err != nil {
		return model.Person{}, fmt.Errorf("Cannot retrieve person from db: %w", err)
	}

	return person, nil
}

func (db *RDSdb) GetAllPersons() ([]model.Person, error) {
	var persons []model.Person

	keys, err := db.Client.Do("KEYS", "person:*").Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve persons from db: %w", err)
	}

	for _, k := range keys.([]interface{}) {
		var person model.Person
		reply, err := db.Client.Do("GET", k.([]byte)).Result()
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve persons from db: %w", err)
		}

		if err := json.Unmarshal(reply.([]byte), &person); err != nil {
			return nil, fmt.Errorf("Failed to retrieve persons from db: %w", err)
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (db *RDSdb) UpdatePerson(id int, p model.Person) (model.Person, error) {
	person, err := json.Marshal(p)
	if err != nil {
		return model.Person{}, fmt.Errorf("Cannot update person in db: %w", err)
	}

	_, err = db.Client.Do("SET", "person:"+strconv.Itoa(p.Id), person).Result()
	if err != nil {
		return model.Person{}, fmt.Errorf("Cannot update person in db: %w", err)
	}

	return p, nil
}

func (db *RDSdb) DeletePerson(id int) error {
	_, err := db.Client.Do("DEL", "person:"+strconv.Itoa(id)).Result()
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}
