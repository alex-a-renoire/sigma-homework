package grpccontroller

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
)

type GRPCСontroller struct {
	remoteStorage pb.StorageServiceClient
}

func New(db pb.StorageServiceClient) GRPCСontroller {
	return GRPCСontroller{
		remoteStorage: db,
	}
}

func (s GRPCСontroller) AddPerson(p model.Person) (uuid.UUID, error) {
	resp, err := s.remoteStorage.AddPerson(context.Background(), &pb.AddPersonRequest{
		Name: p.Name,
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add person: %w", err)
	}

	id, err := uuid.Parse(resp.Value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert types protobuf to postgres: %w", err)
	}

	return id, nil
}

func (s GRPCСontroller) GetPerson(id uuid.UUID) (model.Person, error) {
	p, err := s.remoteStorage.GetPerson(context.Background(), &pb.UUID{

		Value: id.String(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return model.Person{}, fmt.Errorf("failed to get person, failed to parse status or not a grpc error type: %w", err)
		}

		if st.Code() == codes.NotFound {
			return model.Person{}, model.ErrNotFound
		}

		return model.Person{}, fmt.Errorf("failed to get person: %w", err)
	}

	id, err = uuid.Parse(p.Id.Value)
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to convert types protobuf to postgres: %w", err)
	}

	return model.Person{
		Id:   id,
		Name: p.Name,
	}, nil
}

func (s GRPCСontroller) GetAllPersons() ([]model.Person, error) {
	resp, err := s.remoteStorage.GetAllPersons(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch persons: %w", err)
	}

	persons := []model.Person{}

	for _, p := range resp.AllPersons {
		id, err := uuid.Parse(p.Id.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert types protobuf to postgres: %w", err)
		}

		persons = append(persons, model.Person{
			Id:   id,
			Name: p.Name,
		})
	}

	return persons, nil
}

func (s GRPCСontroller) UpdatePerson(id uuid.UUID, person model.Person) error {
	//Check if there is such a person
	_, err := s.remoteStorage.GetPerson(context.Background(), &pb.UUID{
		Value: id.String(),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to get person, failed to parse status or not a grpc error type: %w", err)
		}

		if st.Code() == codes.NotFound {
			return model.ErrNotFound
		}

		return fmt.Errorf("failed to get person: %w", err)
	}

	_, err = s.remoteStorage.UpdatePerson(context.Background(), &pb.Person{
		Id:   &pb.UUID{Value: id.String()},
		Name: person.Name,
	},
	)

	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	return nil
}

func (s GRPCСontroller) DeletePerson(id uuid.UUID) error {
	//Check if there is such a person
	_, err := s.remoteStorage.GetPerson(context.Background(), &pb.UUID{
		Value: id.String(),
	})

	//we assume error is sql.no rows
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to get person, failed to parse status or not a grpc error type: %w", err)
		}

		if st.Code() == codes.NotFound {
			return model.ErrNotFound
		}

		return fmt.Errorf("failed to get person: %w", err)
	}

	_, err = s.remoteStorage.DeletePerson(context.Background(), &pb.DeletePersonRequest{
		Id: &pb.UUID{Value: id.String()},
	})

	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	return nil
}

///////
//CSV//
///////

func (s GRPCСontroller) ProcessCSV(file multipart.File) error {
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
		if _, err = s.remoteStorage.GetPerson(context.Background(), &pb.UUID{
			Value: id.String(),
		}); err == nil {
			if err = s.UpdatePerson(p.Id, p); err != nil {
				return fmt.Errorf("failed to update person in db: %w", err)
			}
			continue
		}

		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("failed to get person, failed to parse status or not a grpc error type: %w", err)
		}

		if st.Code() == codes.NotFound {
			if _, err = s.AddPerson(p); err != nil {
				return fmt.Errorf("failed to add person to db: %w", err)
			}
			continue
		}

		return fmt.Errorf("failed to process csv: %w", err)
	}
}

func (s GRPCСontroller) DownloadPersonsCSV() ([]byte, error) {
	//Ask the service to process action
	persons, err := s.GetAllPersons()
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
