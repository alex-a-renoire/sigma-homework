package service

import (
	"context"
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

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

func (s GRPCPersonService) UpdatePerson(id int, person model.Person) (model.Person, error) {
	//Check if there is such a person
	resp, err := s.remoteStorage.GetPerson(context.Background(), &pb.GetPersonRequest{
		Id: int32(id),
	})

	//we assume error is sql.no rows
	if err != nil {
		return model.Person{}, fmt.Errorf("there is no such person: %w", err)
	}

	resp, err = s.remoteStorage.UpdatePerson(context.Background(), &pb.UpdatePersonRequest{
		Id:   int32(id),
		Name: person.Name,
	})

	if err != nil {
		return model.Person{}, fmt.Errorf("failed to update person: %w", err)
	}

	return model.Person{
		Id:   int(resp.Id),
		Name: resp.Name,
	}, nil
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
