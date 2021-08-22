package service

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jszwec/csvutil"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
)

//TODO validation
//TODO
type GRPCPersonService struct {
	remoteStorage pb.StorageServiceClient
}

func NewGRPC(db pb.StorageServiceClient) GRPCPersonService {
	return GRPCPersonService{
		remoteStorage: db,
	}
}

func (s GRPCPersonService) AddPerson(name string) (int, error) {
	resp, err := s.remoteStorage.AddPerson(context.Background(), &pb.AddPersonRequest{
		Name: name,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to add person: %w", err)
	}

	return int(resp.Id), nil
}

func (s GRPCPersonService) GetPerson(id int) (model.Person, error) {
	resp, err := s.remoteStorage.GetPerson(context.Background(), &pb.GetPersonRequest{
		Id: int32(id),
	})
	if err != nil {
		if errors.Is(err, redis.Nil) || errors.Is(err, sql.ErrNoRows) {
			return model.Person{}, fmt.Errorf("no such record")
		}
		return model.Person{}, fmt.Errorf("failed to get person: %w", err)
	}

	return model.Person{
		Id:   int(resp.Id),
		Name: resp.Name,
	}, nil
}

func (s GRPCPersonService) GetAllPersons() ([]model.Person, error) {
	resp, err := s.remoteStorage.GetAllPersons(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch persons: %w", err)
	}

	persons := []model.Person{}

	for _, p := range resp.AllPersons {
		persons = append(persons, model.Person{
			Id:   int(p.Id),
			Name: p.Name,
		})
	}

	return persons, nil
}

func (s GRPCPersonService) UpdatePerson(id int, person model.Person) error {
	//Check if there is such a person
	_, err := s.remoteStorage.GetPerson(context.Background(), &pb.GetPersonRequest{
		Id: int32(id),
	})

	if err != nil {
		if errors.Is(err, redis.Nil) || errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no such record")
		}
		return fmt.Errorf("failed to get person: %w", err)
	}

	log.Printf("Id service %d", int32(id))
	_, err = s.remoteStorage.UpdatePerson(context.Background(), &pb.Person{
		Id:   int32(id),
		Name: person.Name,
	},
	)

	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	return nil
}

func (s GRPCPersonService) DeletePerson(id int) error {
	//Check if there is such a person
	_, err := s.remoteStorage.GetPerson(context.Background(), &pb.GetPersonRequest{
		Id: int32(id),
	})

	//we assume error is sql.no rows
	if err != nil {
		return fmt.Errorf("there is no such person: %w", err)
	}

	_, err = s.remoteStorage.DeletePerson(context.Background(), &pb.DeletePersonRequest{
		Id: int32(id),
	})

	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	return nil
}

///////
//CSV//
///////

func (s GRPCPersonService) ProcessCSV(file multipart.File) error {
	//Parse CSV
	reader := csv.NewReader(file)
	reader.Read()

	//loop of reading
	for i := 0; ; i++ {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF { //TODO: replace with errors.Is
				if i == 0 {
					//if there's only headers and no values
					return fmt.Errorf("Malformed csv file: %w", err)
				} else {
					//end of the file
					return nil
				}
			} else {
				return fmt.Errorf("Error reading file: %w", err)
			}
		}

		//malformed csv handling
		if len(record) != 2 {
			return fmt.Errorf("Malformed csv file: wrong number of fields")
		}
		if record[0] == "" || record[1] == "" {
			return fmt.Errorf("malformed csv file")
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return fmt.Errorf("malformed id, should be a number: %w", err)
		}

		p := model.Person{
			Id:   id,
			Name: record[1],
		}

		//handle situation when there is such a record and we are updating
		_, err = s.GetPerson(p.Id)
		if err == nil {
			_, err := s.UpdatePerson(p.Id, p.Name)
			if err != nil {
				return fmt.Errorf("failed to update person in db: %w", err)
			}
			return nil
		}

		//If person is not in db, add it with a new id
		if err != sql.ErrNoRows && err != redis.Nil {
			return fmt.Errorf("failed to get person: %w", err)
		}

		if _, err = s.AddPerson(p.Name); err != nil {
			return fmt.Errorf("failed to add person to db: %w", err)
		}
	}
}

func (s GRPCPersonService) DownloadPersonsCSV() ([]byte, error) {
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
