package pgstorage

import (
	"database/sql"
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
)

type PGPersonStorage struct {
	db *sql.DB
}

func New(addr string) (*PGPersonStorage, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &PGPersonStorage{
		db: db,
	}, nil
}

func (s *PGPersonStorage) AddPerson(p model.Person) (int, error) {
	row := d.db.QueryRow("INSERT INTO persons(name) VALUES ($1) RETURNING id", p.Name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PGPersonStorage) GetPerson(id int) (model.Person, error) {

}

func (s *PGPersonStorage) GetAllPersons() ([]model.Person, error) {

}

func (s *PGPersonStorage) UpdatePerson(id int, p model.Person) (model.Person, error) {

}

func (s *PGPersonStorage) DeletePerson(id int) error {

}
