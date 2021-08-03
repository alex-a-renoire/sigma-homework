package service

import (
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
)

type ServicePersonStorage interface {
	AddPerson(name string) (int, error)
	GetPerson(id int) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id int, name string) (model.Person, error)
	DeletePerson(id int) error
}

type PersonService struct {
	s storage.Storage
	FunctionMap map[string]interface{}
}

func New(s ServicePersonStorage) PersonService {
	m := map[string]interface{}{
        "AddPerson": s.AddPerson,
        "GetPerson": s.GetPerson,
		"GetAllPersons": s.GetAllPersons,
		"UpdatePerson": s.UpdatePerson,
		"DeletePerson": s.DeletePerson,
	}

	return PersonService{
		s: s,
		FunctionMap: m,
	}
}

//Controller - storage - map - action
//https://habr.com/ru/post/529086/
//https://medium.com/@matryer/golang-advent-calendar-day-five-routing-restful-controllers-edb74e7d4101
//https://gophersaurus.github.io/docs/v1/controllers/

func (ps PersonService) ProcessAction(action model.Action) (string, error) {
	var response string
	var err error = nil
	person := action.Parameters

	if err := person.Validate(); err != nil && action.FuncName != "GetAllPersons" {
		return "Person name or id not specified", err
	}

	switch action.FuncName {
	case "AddPerson":
		id, err := ps.FunctionMap["AddPerson"].(func(string)(int, error))(person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d and name %s added \n", id, person.Name)
		}
	case "GetPerson":
		p, err := ps.FunctionMap["GetPerson"].(func(int) (model.Person, error))(person.Id)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d has name %s \n", p.Id, p.Name)
		}
	case "GetAllPersons":
		p, err := ps.FunctionMap["GetAllPerson"].(func()([]model.Person, error))()
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("All persons in the storage are %v \n", p)
		}
	case "UpdatePerson":
		p, err := ps.FunctionMap["UpdatePerson"].(func(int, string) (model.Person, error))(person.Id, person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d updated with name %s \n", p.Id, p.Name)
		}

	case "DeletePerson":
		if err := ps.FunctionMap["DeletePerson"].(func(int)(error))(person.Id); err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d deleted \n", person.Id)
		}
	default:
		response = fmt.Sprintf("%s is not a valid command. Try again... \n", action.FuncName)
	}

	return response, err
}
