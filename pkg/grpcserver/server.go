package grpcserver

import (
	"context"
	"fmt"
	"log"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

//Declare the GRPC server
type StorageServer struct {
	pb.UnimplementedStorageServiceServer

	DB storage.Storage
}

func NewGRPC(db storage.Storage) *StorageServer {
	return &StorageServer{
		DB: db,
	}
}

func (ss *StorageServer) AddPerson(_ context.Context, in *pb.AddPersonRequest) (*pb.AddPersonResponse, error) {
	log.Println("Add person command received...")
	p := model.Person{
		Name: in.Name,
	}

	id, err := ss.DB.AddPerson(p)
	if err != nil {
		return nil, fmt.Errorf("failed to add person: %w", err)
	}

	return &pb.AddPersonResponse{
		Id: &pb.UUID{Value: id.String()},
	}, nil
}

func (ss *StorageServer) GetPerson(_ context.Context, in *pb.GetPersonRequest) (*pb.Person, error) {
	log.Println("Get person command received...")
	id, err := uuid.FromBytes([]byte(in.Id.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to bytes: %w", err)
	}
	p, err := ss.DB.GetPerson(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return &pb.Person{
		Id:   &pb.UUID{Value: p.Id.String()},
		Name: p.Name,
	}, nil
}

func (ss *StorageServer) GetAllPersons(_ context.Context, in *emptypb.Empty) (*pb.AllPersonsResponse, error) {
	log.Println("Get all persons command received...")
	persons, err := ss.DB.GetAllPersons()
	if err != nil {
		return nil, fmt.Errorf("failed to get the list of persons: %w", err)
	}

	pbPersons := []*pb.Person{}

	for _, p := range persons {
		pbPersons = append(pbPersons, &pb.Person{
			Id:   &pb.UUID{Value: p.Id.String()},
			Name: p.Name,
		})
	}

	return &pb.AllPersonsResponse{
		AllPersons: pbPersons,
	}, nil
}

func (ss *StorageServer) UpdatePerson(_ context.Context, in *pb.Person) (*emptypb.Empty, error) {
	person := model.Person{
		Name: in.Name,
	}

	log.Println("Update person command received...")

	id, err := uuid.FromBytes([]byte(in.Id.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to bytes: %w", err)
	}

	err = ss.DB.UpdatePerson(id, person)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("failed to update person: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (ss *StorageServer) DeletePerson(_ context.Context, in *pb.DeletePersonRequest) (*emptypb.Empty, error) {
	log.Println("Delete person command received...")

	id, err := uuid.FromBytes([]byte(in.Id.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to convert id to bytes: %w", err)
	}
	if err := ss.DB.DeletePerson(id); err != nil {
		return &emptypb.Empty{}, fmt.Errorf("failed to delete person: %w", err)
	}
	return &emptypb.Empty{}, nil
}
