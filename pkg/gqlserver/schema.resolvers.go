package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphql-go-example/graph/generated"
	"graphql-go-example/graph/model"
)

func (r *mutationResolver) AddPerson(ctx context.Context, input model.PersonInput) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePerson(ctx context.Context, id string, input model.PersonInput) (*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePerson(ctx context.Context, id string) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListPersons(ctx context.Context) ([]*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPerson(ctx context.Context, id int) (*model.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
