package csvservice

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PersonService interface {
	AddPerson(p model.Person) (uuid.UUID, error)
	GetPerson(id uuid.UUID) (model.Person, error)
	GetAllPersons() ([]model.Person, error)
	UpdatePerson(id uuid.UUID, person model.Person) error
	DeletePerson(id uuid.UUID) error
}

type CsvProcessor struct {
	srv PersonService
}

func New(srv PersonService) *CsvProcessor {
	return &CsvProcessor{
		srv: srv,
	}
}

func (cp CsvProcessor) DownloadPersonsCSV() ([]byte, error) {
	persons, err := cp.srv.GetAllPersons()
	if err != nil {
		return nil, fmt.Errorf("failed to get all persons from db: %w", err)
	}

	//Marshal persons into csv
	ps, err := csvutil.Marshal(persons)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal persons: %w", err)
	}

	return ps, nil
}

func (cp CsvProcessor) ProcessCSV(file multipart.File) error {
	//Parse CSV
	reader := csv.NewReader(file)
	reader.Read()

	//loop of reading
	for i := 0; ; i++ {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("Error reading file: %w", err)
			}
			if i == 0 {
				return fmt.Errorf("Malformed csv file: there's only headers and no values")
			} else {
				log.Print("end of file")
				//end of the file
				return nil
			}
		}

		//malformed csv handling
		if len(record) != 2 {
			return fmt.Errorf("Malformed csv file: wrong number of fields")
		}
		if record[0] == "" || record[1] == "" {
			return fmt.Errorf("malformed csv file: empty fields")
		}
		id, err := uuid.Parse(record[0])
		if err != nil {
			return fmt.Errorf("malformed id, should be a uuid: %w", err)
		}

		p := model.Person{
			Id:   id,
			Name: record[1],
		}

		//handle situation when there is such a record and we are updating
		if _, err = cp.srv.GetPerson(id); err == nil {
			if err = cp.srv.UpdatePerson(p.Id, p); err != nil {
				return fmt.Errorf("failed to update person in db: %w", err)
			}
			continue
		}

		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to get person, failed to parse status or not a grpc error type: %w", err)
		}

		if st.Code() == codes.NotFound {
			if _, err = cp.srv.AddPerson(p); err != nil {
				return fmt.Errorf("failed to add person to db: %w", err)
			}
			continue
		}

		return fmt.Errorf("failed to process csv: %w", err)
	}
}
