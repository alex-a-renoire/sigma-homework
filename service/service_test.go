package service

import (
	"fmt"
	"testing"

	"github.com/alex-a-renoire/tcp/model"
	"github.com/alex-a-renoire/tcp/pkg/storage"
)

func TestProcessAction(t *testing.T) {
	type args struct {
		s      storage.Storage
		action model.Action
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "AddPerson OK",
			args: args{
				s: storage.MockStorage{
					MockAddPerson: func(_ string) (int, error) {
						return 1, nil
					},
				},
				action: model.Action{
					FuncName: "AddPerson",
					Parameters: model.Person{
						Name: "Bob",
					},
				},
			},
			want: "Person with id 1 and name Bob added \n",
		},
		{
			name: "AddPerson not OK",
			args: args{
				s: storage.MockStorage{
					MockAddPerson: func(_ string) (int, error) {
						return 0, fmt.Errorf("failed to add person to db")
					},
				},
				action: model.Action{
					FuncName: "AddPerson",
					Parameters: model.Person{
						Name: "Bob",
					},
				},
			},
			want: "error: failed to add person to db \n",
		},
		{
			name: "GetPerson OK",
			args: args{
				s: storage.MockStorage{
					MockGetPerson: func(id int) (model.Person, error) {
						return model.Person{
							Id:   id,
							Name: "Bob",
						}, nil
					},
				},
				action: model.Action{
					FuncName: "GetPerson",
					Parameters: model.Person{
						Id: 1,
					},
				},
			},
			want: "Person with id 1 has name Bob \n",
		},
		{
			name: "GetPerson not OK",
			args: args{
				s: storage.MockStorage{
					MockGetPerson: func(id int) (model.Person, error) {
						return model.Person{}, fmt.Errorf("person not found")
					},
				},
				action: model.Action{
					FuncName: "GetPerson",
					Parameters: model.Person{
						Id: 1,
					},
				},
			},
			want: "error: person not found \n",
		},
		{
			name: "UpdatePerson OK",
			args: args{
				s: storage.MockStorage{
					MockUpdatePerson: func(id int, name string) (model.Person, error) {
						return model.Person{
							Id:   id,
							Name: name,
						}, nil
					},
				},
				action: model.Action{
					FuncName: "UpdatePerson",
					Parameters: model.Person{
						Id:   1,
						Name: "Alice",
					},
				},
			},
			want: "Person with id 1 updated with name Alice \n",
		},
		{
			name: "UpdatePerson not OK",
			args: args{
				s: storage.MockStorage{
					MockUpdatePerson: func(id int, name string) (model.Person, error) {
						return model.Person{}, fmt.Errorf("person not found")
					},
				},
				action: model.Action{
					FuncName: "UpdatePerson",
					Parameters: model.Person{
						Id:   1,
						Name: "Alice",
					},
				},
			},
			want: "error: person not found \n",
		},
		{
			name: "DeletePerson OK",
			args: args{
				s: storage.MockStorage{
					MockDeletePerson: func(_ int) error {
						return nil
					},
				},
				action: model.Action{
					FuncName: "DeletePerson",
					Parameters: model.Person{
						Id: 1,
					},
				},
			},
			want: "Person with id 1 deleted \n",
		},
		{
			name: "DeletePerson not OK",
			args: args{
				s: storage.MockStorage{
					MockDeletePerson: func(_ int) error {
						return fmt.Errorf("person not found")
					},
				},
				action: model.Action{
					FuncName: "DeletePerson",
					Parameters: model.Person{
						Id: 1,
					},
				},
			},
			want: "error: person not found \n",
		},
		{
			name: "Invalid action",
			args: args{
				s: nil,
				action: model.Action{
					FuncName:   "InvalidAction",
					Parameters: model.Person{},
				},
			},
			want: "InvalidAction is not a valid command. Try again... \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ProcessAction(tt.args.s, tt.args.action); got != tt.want {
				t.Errorf("ProcessAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
