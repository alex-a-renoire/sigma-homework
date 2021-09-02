package csvservice

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/alex-a-renoire/sigma-homework/service/authservice"
	"github.com/alex-a-renoire/sigma-homework/service/personservice"
	"github.com/google/uuid"
)

func TestCsvProcessor_ProcessCSV(t *testing.T) {
	type fields struct {
		db storage.MockStorage
	}
	tests := []struct {
		name    string
		file    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Empty file",
			file:    "testdata/empty.csv",
			wantErr: true,
		},
		{
			name:    "Only headers",
			file:    "testdata/onlyheaders.csv",
			wantErr: true,
		},
		{
			name:    "Wrong number of fields",
			file:    "testdata/toomanyfields.csv",
			wantErr: true,
		},
		{
			name:    "Empty fields",
			file:    "testdata/emptyfields.csv",
			wantErr: true,
		},
		{
			name:    "Wrong ID format",
			file:    "testdata/wrongid.csv",
			wantErr: true,
		},
		{
			name: "Add person OK",
			file: "testdata/add.csv",
			fields: fields{
				db: storage.MockStorage{
					MockGetPerson: func(id uuid.UUID) (model.Person, error) {
						return model.Person{}, model.ErrNotFound
					},
					MockAddPerson: func(_ model.Person) (uuid.UUID, error) {
						return uuid.New(), nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Update person OK",
			file: "testdata/rename.csv",
			fields: fields{
				db: storage.MockStorage{
					MockGetPerson: func(id uuid.UUID) (model.Person, error) {
						return model.Person{
							Id:   id,
							Name: "John",
						}, nil
					},
					MockUpdatePerson: func(_ uuid.UUID, p model.Person) error {
						return nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		auth := authservice.New("some_weak_secret")
		srv := personservice.New(tt.fields.db, auth)

		file, _ := os.Open(tt.file)
		reader := csv.NewReader(file)

		t.Run(tt.name, func(t *testing.T) {
			cp := CsvProcessor{
				srv: srv,
			}
			if err := cp.ProcessCSV(*reader); (err != nil) != tt.wantErr {
				t.Errorf("CsvProcessor.ProcessCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
