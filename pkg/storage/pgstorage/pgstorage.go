package pgstorage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

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
	row := s.db.QueryRow("INSERT INTO persons(name) VALUES ($1) RETURNING id", p.Name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add person to db")
	}

	return id, nil
}

func (s *PGPersonStorage) GetPerson(id int) (model.Person, error) {
	row := s.db.QueryRow("SELECT id, name FROM persons WHERE id=$1", id)

	var p model.Person
	if err := row.Scan(&p.Id, &p.Name); err != nil {
		return model.Person{}, fmt.Errorf("failed to get person from db: %w", err)
	}
	return p, nil
}

func (s *PGPersonStorage) GetAllPersons() ([]model.Person, error) {
	rows, err := s.db.Query("SELECT id, name FROM persons")
	if err != nil {
		return nil, fmt.Errorf("failed to get persons from db: %w", err)
	}
	defer rows.Close()

	slicePersons := []model.Person{}
	for rows.Next() {
		var u model.Person
		if err := rows.Scan(
			&u.Id,
			&u.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to get persons from db: %w", err)
		}
		slicePersons = append(slicePersons, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get persons from db: %w", err)
	}
	return slicePersons, nil
}

func (s *PGPersonStorage) UpdatePerson(id int, p model.Person) error {
	_, err := s.db.Exec("UPDATE persons SET name = $1 WHERE id=$2", p.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	return nil
}

func (s *PGPersonStorage) DeletePerson(id int) error {
	_, err := s.db.Exec("DELETE FROM persons WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}
