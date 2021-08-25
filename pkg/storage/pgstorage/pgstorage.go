package pgstorage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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

func (s *PGPersonStorage) AddPerson(p model.Person) (uuid.UUID, error) {
	id := uuid.New()
	if _, err := s.db.Exec("INSERT INTO persons(id, name) VALUES ($1, $2)", id, p.Name); err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to add person to db: %w", err)
	}

	return id, nil
}

func (s *PGPersonStorage) GetPerson(id uuid.UUID) (model.Person, error) {
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

func (s *PGPersonStorage) UpdatePerson(id uuid.UUID, p model.Person) error {
	_, err := s.db.Exec("UPDATE persons SET name = $1 WHERE id=$2", p.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	return nil
}

func (s *PGPersonStorage) DeletePerson(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM persons WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}
