package personservice

import (
	"fmt"
	"testing"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/google/uuid"
)

var id = uuid.New()

func TestPersonService_AddPerson(t *testing.T) {
	type fields struct {
		db PersonStorage
	}
	type args struct {
		p model.AddUpdatePerson
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "AddPerson OK",
			fields: fields{
				db: storage.MockStorage{
					MockAddPerson: func(_ model.Person) (uuid.UUID, error) {
						return uuid.New(), nil
					},
				},
			},
			args: args{
				p: model.AddUpdatePerson{
					Name: "Boris",
				},
			},
			wantErr: false,
		},
		{
			name: "AddPerson !OK",
			fields: fields{
				db: storage.MockStorage{
					MockAddPerson: func(_ model.Person) (uuid.UUID, error) {
						return uuid.Nil, fmt.Errorf("weird error")
					},
				},
			},
			args: args{
				p: model.AddUpdatePerson{
					Name: "Boris",
				},
			},
			wantErr: true,
		},
		{
			name: "AddPerson ID not empty",
			args: args{
				p: model.AddUpdatePerson{
					Id:   id,
					Name: "Boris",
				},
			},
			wantErr: true,
		},
		{
			name: "AddPerson Name empty",
			args: args{
				p: model.AddUpdatePerson{
					Id: id,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PersonService{
				db: tt.fields.db,
			}

			_, err := s.AddPerson(tt.args.p)
			if err != nil && !tt.wantErr {
				t.Errorf("PersonService.AddPerson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_GetPerson(t *testing.T) {
	type fields struct {
		db PersonStorage
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "GetPerson ID nil",
			args: args{
				id: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "GetPerson OK",
			fields: fields{
				db: storage.MockStorage{
					MockGetPerson: func(_ uuid.UUID) (model.Person, error) {
						return model.Person{}, nil
					},
				},
			},
			args: args{
				id: id,
			},
			wantErr: false,
		},
		{
			name: "GetPerson !OK",
			fields: fields{
				db: storage.MockStorage{
					MockGetPerson: func(_ uuid.UUID) (model.Person, error) {
						return model.Person{}, fmt.Errorf("weird DB error")
					},
				},
			},
			args: args{
				id: id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PersonService{
				db: tt.fields.db,
			}
			_, err := s.GetPerson(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonService.GetPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPersonService_DeletePerson(t *testing.T) {
	type fields struct {
		db PersonStorage
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "DeletePerson ID nil",
			args: args{
				id: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "DeletePerson OK",
			fields: fields{
				db: storage.MockStorage{
					MockDeletePerson: func(_ uuid.UUID) error {
						return nil
					},
				},
			},
			args: args{
				id: id,
			},
			wantErr: false,
		},
		{
			name: "DeletePerson !OK",
			fields: fields{
				db: storage.MockStorage{
					MockDeletePerson: func(_ uuid.UUID) error {
						return fmt.Errorf("weird DB error")
					},
				},
			},
			args: args{
				id: id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PersonService{
				db: tt.fields.db,
			}
			err := s.DeletePerson(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonService.GetPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPersonService_UpdatePerson(t *testing.T) {
	type fields struct {
		db PersonStorage
	}
	type args struct {
		id uuid.UUID
		p  model.AddUpdatePerson
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdatePerson OK",
			fields: fields{
				db: storage.MockStorage{
					MockUpdatePerson: func(id uuid.UUID, _ model.Person) error {
						return nil
					},
				},
			},
			args: args{
				id: id,
				p: model.AddUpdatePerson{
					Name: "Boris",
				},
			},
			wantErr: false,
		},
		{
			name: "UpdatePerson !OK",
			fields: fields{
				db: storage.MockStorage{
					MockUpdatePerson: func(id uuid.UUID, _ model.Person) error {
						return fmt.Errorf("some weird DB error")
					},
				},
			},
			args: args{
				id: id,
				p: model.AddUpdatePerson{
					Name: "Boris",
				},
			},
			wantErr: true,
		},
		{
			name: "UpdatePerson ID nil",
			args: args{
				id: uuid.Nil,
				p: model.AddUpdatePerson{
					Id:   id,
					Name: "Boris",
				},
			},
			wantErr: true,
		},
		{
			name: "UpdatePerson Name empty",
			args: args{
				p: model.AddUpdatePerson{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PersonService{
				db: tt.fields.db,
			}
			if err := s.UpdatePerson(tt.args.id, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.UpdatePerson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
