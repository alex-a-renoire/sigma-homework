package mongostorage

import (
	"context"
	"fmt"
	"log"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPersonStorage struct {
	Client *mongo.Client
}

func New(addr string, user string, password string) (*MongoPersonStorage, error) {
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //TODO: ask what to do with cancel function?
	log.Print("mongodb:" + addr)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://"+user+":"+password+"@"+addr+"/persons"))
	if err != nil {
		return nil, fmt.Errorf("failed to open the db connection: %w", err)
	}

	//TODO: ask how to close the connection?

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("Conncetion is not yet established: %w", err)
	}
	fmt.Println("Connected to MongoDB!")

	return &MongoPersonStorage{
		Client: client,
	}, nil
}

func (mg *MongoPersonStorage) AddPerson(p model.Person) (uuid.UUID, error) {
	collection := mg.Client.Database("persons").Collection("persons")

	p.Id = uuid.New()
	_, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to add person to db: %w", err)
	}

	return p.Id, nil
}

func (mg *MongoPersonStorage) GetPerson(id uuid.UUID) (model.Person, error) {
	collection := mg.Client.Database("persons").Collection("persons")
	var p model.Person
	if err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&p); err != nil {
		return model.Person{}, fmt.Errorf("Failed to get person from db: %w", err)
	}

	return p, nil
}

func (mg *MongoPersonStorage) GetAllPersons() ([]model.Person, error) {
	return nil, nil
}

func (mg *MongoPersonStorage) UpdatePerson(id uuid.UUID, person model.Person) error {
	collection := mg.Client.Database("persons").Collection("persons")

	_id, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	filter := bson.D{{Key: "id", Value: _id}}
	update := bson.D{{"$set", primitive.E{Key: "name", Value: person.Name}}}

	collection.FindOneAndUpdate(context.TODO(), filter, update)

	return nil
	//TODO find how to check if updates were successfull
}

func (mg *MongoPersonStorage) DeletePerson(id uuid.UUID) error {
	collection := mg.Client.Database("persons").Collection("persons")

	_id, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	opts := options.Delete().SetCollation(&options.Collation{})
	if _, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts); err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}
