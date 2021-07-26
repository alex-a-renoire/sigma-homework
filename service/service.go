package service

import (
	"fmt"

	"github.com/alex-a-renoire/tcp/model"
	"github.com/alex-a-renoire/tcp/pkg/storage"
)

//TODO детали ошибки должны прилетать из стораджа
func ProcessAction(s storage.Storage, action model.Action) string {
	var response string
	person := action.Parameters

	switch action.FuncName {
	case "AddPerson":
		id, err := s.AddPerson(person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d and name %s added \n", id, person.Name)
		}
	case "UpdatePerson":
		p, err := s.UpdatePerson(person.Id, person.Name)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d updated with name %s \n", p.Id, p.Name)
		}
	case "GetPerson":
		p, err := s.GetPerson(person.Id)
		if err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d has name %s \n", p.Id, p.Name)
		}
	case "DeletePerson":
		if err := s.DeletePerson(person.Id); err != nil {
			response = fmt.Sprintf("error: %s \n", err)
		} else {
			response = fmt.Sprintf("Person with id %d deleted \n", person.Id)
		}
	default:
		response = fmt.Sprintf("%s is not a valid command. Try again... \n", action.FuncName)
	}

	return response
}
