package controllers

import (
	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service"
)

type PersonController struct {
	service service.PersonService
}

func New(ps service.PersonService) PersonController {
	return PersonController{
		service: ps,
	}
}

func (pc PersonController) AddPerson(item model.Person) (string, error) {
	return pc.service.ProcessAction(model.Action{
		FuncName:   "AddPerson",
		Parameters: item,
	},
	)
}

func (pc PersonController) GetPerson(id int) (string, error) {
	return pc.service.ProcessAction(model.Action{
		FuncName:   "GetPerson",
		Parameters: model.Person{Id: id},
	},
	)
}

func (pc PersonController) GetAllPersons() (string, error) {
	return pc.service.ProcessAction(model.Action{
		FuncName:   "GetAllPersons",
		Parameters: model.Person{},
	},
	)
}

func (pc PersonController) UpdatePerson(item model.Person) (string, error) {
	return pc.service.ProcessAction(model.Action{
		FuncName:   "UpdatePerson",
		Parameters: item,
	},
	)
}

func (pc PersonController) DeletePerson(id int) (string, error) {
	return pc.service.ProcessAction(model.Action{
		FuncName: "DeletePerson",
		Parameters: model.Person{
			Id: id,
		},
	},
	)
}