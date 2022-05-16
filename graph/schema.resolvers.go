package graph

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aniruddha2000/hackernews/graph/generated"
	"github.com/aniruddha2000/hackernews/graph/model"
	"github.com/aniruddha2000/hackernews/internal/auth"
	"github.com/aniruddha2000/hackernews/internal/links"
	"github.com/aniruddha2000/hackernews/internal/users"
	"github.com/aniruddha2000/hackernews/pkg/jwt"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}

	newLink := links.Link{
		Title:   input.Title,
		Address: input.Address,
		User:    user,
	}
	linkID := newLink.Save()

	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   newLink.Title,
		Address: newLink.Address,
		User:    graphqlUser,
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user = users.User{
		Username: input.Username,
		Password: input.Password,
	}
	err := user.Create()
	if err != nil {
		return "", err
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) UpdateLink(ctx context.Context, id string, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, errors.New("access denied")
	}

	link := links.Link{
		Title:   input.Title,
		Address: input.Address,
		User:    user,
	}
	rowsAffected, err := link.Update(id)
	if err != nil {
		return &model.Link{}, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("zero rows affected")
	}

	return &model.Link{
		ID:      id,
		Title:   link.Title,
		Address: link.Address,
		User: &model.User{
			ID:   user.ID,
			Name: user.Username,
		},
	}, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, id string) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", errors.New("access denied")
	}

	rowsAffected, err := links.Delete(id, user.Username)
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", errors.New("zero rows affected")
	}
	return fmt.Sprintf("%v rows deleted", rowsAffected), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user = users.User{
		Username: input.Username,
		Password: input.Password,
	}

	correct := user.Authenticate()
	if !correct {
		return "", errors.New("wrong username or password")
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var result []*model.Link

	allLinks := links.GetAll()
	for _, link := range allLinks {
		graphqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		result = append(result, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: graphqlUser})
	}
	return result, nil
}

func (r *queryResolver) Link(ctx context.Context, id string) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, errors.New("access denied")
	}

	link, err := links.Get(id, user.Username)
	if err != nil {
		return &model.Link{}, err
	}
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
