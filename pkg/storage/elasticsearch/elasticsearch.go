package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
)

const (
	index = "persons"
)

type SearchResult struct {
	Source model.Person `json:"_source,omitempty"`
}

type ElasticPersonStorage struct {
	client elasticsearch.Client
}

func New(address string) (ElasticPersonStorage, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{address},
	})

	if err != nil {
		return ElasticPersonStorage{}, fmt.Errorf("failed to create client: %w", err)
	}

	return ElasticPersonStorage{
		client: *client,
	}, nil
}

func (eps ElasticPersonStorage) AddPerson(p model.Person) (uuid.UUID, error) {
	p.Id = uuid.New()

	person, err := json.Marshal(p)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal person: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      index,
		DocumentID: p.Id.String(),
		Body:       bytes.NewReader(person),
	}

	res, err := req.Do(context.Background(), &eps.client)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add person to db (elastic): %w", err)
	}
	defer res.Body.Close()

	return p.Id, nil
}

func (eps ElasticPersonStorage) GetPerson(id uuid.UUID) (model.Person, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id.String(),
	}

	res, err := req.Do(context.Background(), &eps.client)
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to get person: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.Person{}, model.ErrNotFound
	}

	if res.IsError() {
		return model.Person{}, fmt.Errorf("failed to get person: %s", res.String())
	}

	log.Printf("Res-body: %s", res.Body)

	sr := SearchResult{}
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return model.Person{}, fmt.Errorf("failed to unmarshal person: %s", res.Body)
	}

	return sr.Source, nil
}

func (eps ElasticPersonStorage) GetAllPersons() ([]model.Person, error) {
	req := esapi.SearchRequest{
		Index: []string{index},
	}

	resp, err := req.Do(context.Background(), &eps.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get persons from db: %w", err)
	}

	var persons []model.Person
	err = json.NewDecoder(resp.Body).Decode(&persons)
	if err != nil {
		return nil, fmt.Errorf("sdfsdf: %w", err)
	}

	return persons, nil
}

func (eps ElasticPersonStorage) UpdatePerson(id uuid.UUID, person model.Person) error {
	p, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("failed to marshal person: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id.String(),
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, p))),
	}

	res, err := req.Do(context.Background(), &eps.client)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("failed to update person: %s", res.String())
	}

	return nil
}

func (eps ElasticPersonStorage) DeletePerson(id uuid.UUID) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id.String(),
	}

	res, err := req.Do(context.Background(), &eps.client)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("failed to delete person: %s", res.String())
	}
	return nil
}
