package dummytcp

type Action struct {
	FuncName string `json:"func_name,omitempty"`
	Parameters     Person `json:"data,omitempty"`
}
