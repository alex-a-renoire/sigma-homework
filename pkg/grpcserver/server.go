package grpcserver

import (
	"context"
	"fmt"
	"log"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
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
		Id: int32(id),
	}, nil
}

func (ss *StorageServer) GetPerson(_ context.Context, in *pb.GetPersonRequest) (*pb.Person, error) {
	log.Println("Get person command received...")
	p, err := ss.DB.GetPerson(int(in.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return &pb.Person{
		Id:   int32(p.Id),
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
			Id:   int32(p.Id),
			Name: p.Name,
		})
	}

	return &pb.AllPersonsResponse{
		AllPersons: pbPersons,
	}, nil
}

func (ss *StorageServer) UpdatePerson(_ context.Context, in *pb.UpdatePersonRequest) (*pb.Person, error) {
	person := model.Person{
		Name: in.Name,
	}

	log.Println("Update person command received...")
	p, err := ss.DB.UpdatePerson(int(in.Id), person)
	if err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	return &pb.Person{
		Id:   int32(p.Id),
		Name: p.Name,
	}, nil
}

func (ss *StorageServer) DeletePerson(_ context.Context, in *pb.DeletePersonRequest) (*emptypb.Empty, error) {
	log.Println("Delete person command received...")
	if err := ss.DB.DeletePerson(int(in.Id)); err != nil {
		return &emptypb.Empty{}, fmt.Errorf("failed to delete person: %w", err)
	}
	return &emptypb.Empty{}, nil
}
