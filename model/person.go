package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Person struct {
	Id   uuid.UUID `json:"id,omitempty" csv:"id, omitempty" bson:"_id,omitempty"`
	Name string    `json:"name,omitempty" csv:"name" bson:"name,omitempty"`
}

func (p *Person) String() string {
	return fmt.Sprintf("person with id %d and name %s", p.Id, p.Name)
}

func (p *Person) Validate() error {
	if p.Id == uuid.Nil || p.Name == "" {
		return fmt.Errorf("person name or id are not specified. At least one of the two should be specified")
	}
	return nil
}
