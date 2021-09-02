package model

import (
	"fmt"

	"github.com/google/uuid"
)

type PersonAuth struct {
	Id    uuid.UUID `json:"id,omitempty" csv:"id, omitempty" bson:"_id,omitempty"`
	Name  string    `json:"name,omitempty" csv:"name" bson:"name,omitempty"`
	Token string    `json:"token,omitempty"`
}

func (p *PersonAuth) String() string {
	return fmt.Sprintf("person with id %d, name %s and token", p.Id, p.Name)
}

func (p *PersonAuth) Validate() error {
	if p.Id == uuid.Nil || p.Name == "" {
		return fmt.Errorf("person name or id are not specified. At least one of the two should be specified")
	}
	return nil
}
