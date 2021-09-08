package resolver

import "github.com/alex-a-renoire/sigma-homework/service/personservice"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service personservice.PersonService
}

func New(service personservice.PersonService) *Resolver {
	return &Resolver{
		service: service,
	}
}
