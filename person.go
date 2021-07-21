package dummytcp

import "fmt"

type Person struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (p *Person) String() string {
	return fmt.Sprintf("person with id %d and name %s", p.Id, p.Name)
}