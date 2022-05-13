package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aniruddha2000/hackernews/graph/generated"
	"github.com/aniruddha2000/hackernews/graph/model"
	"github.com/aniruddha2000/hackernews/internal/links"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	newLink := links.Links{
		Title:   input.Title,
		Address: input.Address,
	}
	linkID := newLink.Save()

	return &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   newLink.Title,
		Address: newLink.Address,
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateLink(ctx context.Context, id string, input model.NewLink) (*model.Link, error) {
	link := links.Links{
		Title:   input.Title,
		Address: input.Address,
	}
	rowsAffected := link.Update(id)
	if rowsAffected == 0 {
		return nil, errors.New("zero rows affected")
	}
	return &model.Link{
		ID:      id,
		Title:   link.Title,
		Address: link.Address,
	}, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, id string) (string, error) {
	rowsAffected := links.Delete(id)
	if rowsAffected == 0 {
		return "", errors.New("zero rows affected")
	}
	return fmt.Sprintf("%v rows affected", rowsAffected), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var result []*model.Link

	allLinks := links.GetAll()
	for _, link := range allLinks {
		result = append(result, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address})
	}
	return result, nil
}

func (r *queryResolver) Link(ctx context.Context, id string) (*model.Link, error) {
	link := links.Get(id)
	return &model.Link{
		ID:      link.ID,
		Title:   link.Title,
		Address: link.Address,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
