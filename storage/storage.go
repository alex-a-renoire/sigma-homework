package storage

import dummytcp "github.com/alex-a-renoire/tcp"

type Storage interface {
	AddPerson(name string) int
	GetPerson(id int) (dummytcp.Person, error)
	UpdatePerson(id int, name string) (dummytcp.Person, error)
	DeletePerson(id int) error
}
