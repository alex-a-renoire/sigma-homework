package grpccontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
			return model.Person{}, fmt.Errorf("failed to get person in grpc controller, failed to parse status or not a grpc error type: %w", model.ErrNotFound)
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
