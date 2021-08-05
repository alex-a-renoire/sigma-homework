package grpcserver

import (
	"context"
	"fmt"

	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
)

//Declare the GRPC server
type StorageServer struct {
	pb.UnimplementedStorageServiceServer

	DB storage.Storage
}

func (ss *StorageServer) AddPerson(_ context.Context, in *pb.AddPersonRequest) (*pb.AddPersonResponse, error) {
	id, err := ss.DB.AddPerson(in.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to add person: %w", err)
	}

	return &pb.AddPersonResponse{
		Id: int32(id),
	}, nil
}

func (ss *StorageServer) GetPerson(_ context.Context, in *pb.GetPersonRequest) (*pb.Person, error) {
	p, err := ss.DB.GetPerson(int(in.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return &pb.Person{
		Id:   int32(p.Id),
		Name: p.Name,
	}, nil
}

func (ss *StorageServer) GetAllPersons(_ context.Context, in *pb.AllPersonsRequest) (*pb.AllPersonsResponse, error) {
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
	p, err := ss.DB.UpdatePerson(int(in.Id), in.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	return &pb.Person{
		Id:   int32(p.Id),
		Name: p.Name,
	}, nil
}

func (ss *StorageServer) DeletePerson(_ context.Context, in *pb.DeletePersonRequest) (*emptypb.Empty, error) {
	if err := ss.DB.DeletePerson(int(in.Id)); err != nil {
		return &emptypb.Empty{}, fmt.Errorf("failed to delete person: %w", err)
	}
	return &emptypb.Empty{}, nil
}