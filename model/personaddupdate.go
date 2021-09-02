package model

import (
	"fmt"

	"github.com/google/uuid"
)

type AddUpdatePerson struct {
	Id   uuid.UUID `json:"id,omitempty" csv:"id, omitempty" bson:"_id,omitempty"`
	Name string    `json:"name,omitempty" csv:"name" bson:"name,omitempty"`
}

func (p *AddUpdatePerson) String() string {
	return fmt.Sprintf("person with name %s", p.Name)
}

func (p *AddUpdatePerson) Validate() error {
	if p.Id != uuid.Nil {
		return fmt.Errorf("ID should not be specified")
	}
	if p.Name == "" {
		return fmt.Errorf("Name should be specified")
	}
	return nil
}
