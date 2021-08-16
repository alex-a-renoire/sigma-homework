package model

import "fmt"

type Person struct {
	Id   int    `json:"id,omitempty" csv:"id, omitempty"`
	Name string `json:"name,omitempty" csv:"name"`
}

func (p *Person) String() string {
	return fmt.Sprintf("person with id %d and name %s", p.Id, p.Name)
}

func (p *Person) Validate() error {
	if p.Id == 0 && p.Name == "" {
		return fmt.Errorf("person name and id are not specified. At least one of the two should be specified")
	}
	return nil
}
