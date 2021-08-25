package tcpcontroller

import (
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service"
	"github.com/google/uuid"
)

type PersonControllerTCP struct {
	s           service.PersonService
	FunctionMap map[string]interface{}
}

func New(s service.PersonService) PersonControllerTCP {
	m := map[string]interface{}{
		"AddPerson":     s.AddPerson,
		"GetPerson":     s.GetPerson,
		"GetAllPersons": s.GetAllPersons,
		"UpdatePerson":  s.UpdatePerson,
		"DeletePerson":  s.DeletePerson,
	}

	return PersonControllerTCP{
		s:           s,
		FunctionMap: m,
	}
}

func (ps PersonControllerTCP) ProcessAction(action model.Action) (string, error) {
	var response string
	var err error = nil
	person := action.Parameters

	if err := person.Validate(); err != nil && action.FuncName != "GetAllPersons" {
		return "Person name or id not specified", err
	}

	switch action.FuncName {
	case "AddPerson":
		id, err := ps.FunctionMap["AddPerson"].(func(string) (uuid.UUID, error))(person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d and name %s added \n", id, person.Name)
		}
	case "GetPerson":
		p, err := ps.FunctionMap["GetPerson"].(func(uuid.UUID) (model.Person, error))(person.Id)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d has name %s \n", p.Id, p.Name)
		}
	case "GetAllPersons":
		p, err := ps.FunctionMap["GetAllPersons"].(func() ([]model.Person, error))()
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("All persons in the storage are %v \n", p)
		}
	case "UpdatePerson":
		p, err := ps.FunctionMap["UpdatePerson"].(func(uuid.UUID, string) (model.Person, error))(person.Id, person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d updated with name %s \n", p.Id, p.Name)
		}

	case "DeletePerson":
		if err := ps.FunctionMap["DeletePerson"].(func(uuid.UUID) error)(person.Id); err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d deleted \n", person.Id)
		}
	default:
		response = fmt.Sprintf("%s is not a valid command. Try again... \n", action.FuncName)
	}

	return response, err
}
