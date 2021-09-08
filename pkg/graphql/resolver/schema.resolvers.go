package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	graphql1 "github.com/alex-a-renoire/sigma-homework/pkg/graphql"
)

func (r *mutationResolver) AddPerson(ctx context.Context, person graphql1.AddUpdatePerson) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePerson(ctx context.Context, id string, person graphql1.AddUpdatePerson) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePerson(ctx context.Context, id string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllPersons(ctx context.Context) ([]*graphql1.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPerson(ctx context.Context, id string) (*graphql1.Person, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns graphql1.MutationResolver implementation.
func (r *Resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

// Query returns graphql1.QueryResolver implementation.
func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
