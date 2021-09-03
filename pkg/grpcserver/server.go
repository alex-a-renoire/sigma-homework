package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/alex-a-renoire/sigma-homework/model"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageServer struct {
	pb.UnimplementedStorageServiceServer

	DB storage.Storage
}

func NewGRPC(db storage.Storage) *StorageServer {
	return &StorageServer{
		DB: db,
	}
}

func (ss *StorageServer) AddPerson(_ context.Context, in *pb.AddPersonRequest) (*pb.UUID, error) {
	log.Println("Add person command received...")
	p := model.Person{
		Name: in.Name,
	}

	id, err := ss.DB.AddPerson(p)
	if err != nil {
		errStr := fmt.Sprintf("failed to add person: %s", err)
		return nil, status.Error(codes.Internal, errStr)
	}

	return &pb.UUID{
		Value: id.String(),
	}, nil
}

func (ss *StorageServer) GetPerson(_ context.Context, in *pb.UUID) (*pb.Person, error) {
	log.Println("Get person command received...")

	id, err := uuid.Parse(in.Value)
	if err != nil {
		errStr := fmt.Sprintf("failed to parse uuid in grpc server: %s", err)
		return nil, status.Error(codes.Internal, errStr)
	}

	p, err := ss.DB.GetPerson(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, redis.Nil) {
			return nil, status.Error(codes.NotFound, "no rows found")
		}
		errStr := fmt.Sprintf("failed to get person in grpc server: %s", err)
		return nil, status.Error(codes.Internal, errStr)
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
		errStr := fmt.Sprintf("failed to get the list of persons: %s", err)
		return nil, status.Error(codes.Internal, errStr)
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
	log.Println("Update person command received...")

	id, err := uuid.Parse(in.Id.Value)
	if err != nil {
		errStr := fmt.Sprintf("failed to parse uuid: %s", err)
		return nil, status.Error(codes.Internal, errStr)
	}

	person := model.Person{
		Id:   id,
		Name: in.Name,
	}

	err = ss.DB.UpdatePerson(id, person)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, redis.Nil) {
			status.Error(codes.NotFound, "no rows found")
		}

		errStr := fmt.Sprintf("failed to update person: %s", err)
		return &emptypb.Empty{}, status.Error(codes.Internal, errStr)
	}

	return &emptypb.Empty{}, nil
}

func (ss *StorageServer) DeletePerson(_ context.Context, in *pb.DeletePersonRequest) (*emptypb.Empty, error) {
	log.Println("Delete person command received...")

	id, err := uuid.Parse(in.Id.Value)
	if err != nil {
		errStr := fmt.Sprintf("failed to parse uuid: %s", err)
		return &emptypb.Empty{}, status.Error(codes.Internal, errStr)
	}

	if err := ss.DB.DeletePerson(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, redis.Nil) {
			status.Error(codes.NotFound, "no rows found")
		}

		errStr := fmt.Sprintf("failed to delete person: %s", err)
		return &emptypb.Empty{}, status.Error(codes.Internal, errStr)
	}
	return &emptypb.Empty{}, nil
}
