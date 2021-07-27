package model

import "fmt"

type Action struct {
	FuncName   string `json:"func_name"`
	Parameters Person `json:"data"`
}

func (a *Action) Validate() error {
	if a.FuncName == "" {
		return fmt.Errorf("Action is not specified, field func_name is empty or wrongly spelled")
	}

	//validate person
	return a.Parameters.Validate()
}
