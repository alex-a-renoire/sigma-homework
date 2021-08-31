package personservice

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/alex-a-renoire/sigma-homework/model"
// 	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// )

// func TestPersonService_AddPerson(t *testing.T) {
// 	type fields struct {
// 		db PersonStorage
// 	}
// 	type args struct {
// 		p model.AddUpdatePerson
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "AddPerson OK",
// 			fields: fields{
// 				db: storage.MockStorage{
// 					MockAddPerson: func(_ model.Person) (uuid.UUID, error) {
// 						return uuid.New(), nil
// 					},
// 				},
// 			},
// 			args: args{
// 				p: model.AddUpdatePerson{
// 					Name: "Boris",
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		id := uuid.New()
// 		t.Run(tt.name, func(t *testing.T) {
// 			id
// 			assertions := require.New(t)
// 			s := PersonService{
// 				db: tt.fields.db,
// 			}
// 			got, err := s.AddPerson(tt.args.p)
// 			// if (err != nil) != tt.wantErr {
// 			// 	t.Errorf("PersonService.AddPerson() error = %v, wantErr %v", err, tt.wantErr)
// 			// 	return
// 			// }
// 			assertions.Equal(uuid)

// 		})
// 	}
// }

// func TestPersonService_GetPerson(t *testing.T) {
// 	type fields struct {
// 		db PersonStorage
// 	}
// 	type args struct {
// 		id uuid.UUID
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    uuid.UUID
// 		wantErr bool
// 	}{
// 		{
// 			name: "GetPerson OK",
// 			fields: fields{
// 				db: storage.MockStorage{
// 					MockAddPerson: func(_ model.Person) (uuid.UUID, error) {
// 						id, _ := uuid.Parse("")
// 						return id, nil
// 					},
// 				},
// 			},
// 			args: args{
// 				p: model.AddUpdatePerson{
// 					Name: "Boris",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := PersonService{
// 				db: tt.fields.db,
// 			}
// 			got, err := s.GetPerson(tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("PersonService.GetPerson() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("PersonService.GetPerson() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPersonService_UpdatePerson(t *testing.T) {
// 	type fields struct {
// 		db PersonStorage
// 	}
// 	type args struct {
// 		id uuid.UUID
// 		p  model.AddUpdatePerson
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := PersonService{
// 				db: tt.fields.db,
// 			}
// 			if err := s.UpdatePerson(tt.args.id, tt.args.p); (err != nil) != tt.wantErr {
// 				t.Errorf("PersonService.UpdatePerson() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPersonService_DeletePerson(t *testing.T) {
// 	type fields struct {
// 		db PersonStorage
// 	}
// 	type args struct {
// 		id uuid.UUID
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := PersonService{
// 				db: tt.fields.db,
// 			}
// 			if err := s.DeletePerson(tt.args.id); (err != nil) != tt.wantErr {
// 				t.Errorf("PersonService.DeletePerson() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
