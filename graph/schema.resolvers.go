package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.32

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "go-blog-graphql/graph/model"
	dbModel "go-blog-graphql/models"
)

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, authorID string, input model.EditedPost) (*dbModel.Post, error) {
	authorId, err := strconv.ParseUint(authorID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	NewPost := dbModel.Post{
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Author: dbModel.Author {ID: uint(authorId)},
	}

	if err := r.Database.Model(&model.Post{}).Create(&NewPost).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &NewPost, nil
}

// UpdatePost is the resolver for the updatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input model.EditedPost) (*model.Post, error) {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	UpdatePost := dbModel.Post{
		Title:     input.Title,
		Content:   input.Content,
		UpdatedAt: time.Now(),
	}

	if err := r.Database.Model(&model.Post{}).Where("id=?", uint(postId)).Updates(&UpdatePost).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &UpdatePost, nil
}

// AddAuthor is the resolver for the addAuthor field.
func (r *mutationResolver) AddAuthor(ctx context.Context, name string) (*model.Author, error) {
	NewAuthor := r.Database.Model(&model.Author{}).Create(&dbModel.Author{Name: name})
	return NewAuthor, nil
}

// UpdateAuthor is the resolver for the updateAuthor field.
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id string, name string) (*model.Author, error) {
	authorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	UpdatedAuthor := r.Database.Model(&model.Author{}).Where("id = ?", uint(authorId)).Update("name", name)
	return UpdatedAuthor, nil
}

// GetAllPosts is the resolver for the getAllPosts field.
func (r *queryResolver) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts := []*model.Post{}

	GetPosts := r.Database.Model(&posts).Preload("Author").Find(&posts)

	if GetPosts.Error != nil {
		fmt.Println(GetPosts.Error)
		return nil, GetPosts.Error
	}
	return posts, nil
}

// GetPostByID is the resolver for the getPostByID field.
func (r *queryResolver) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	post := model.Post{}

	if err := r.Database.Model(&model.Post{}).Preload("Author").Find(&post, uint(postId)).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &post, nil
}

// GetPostsByAuthorID is the resolver for the getPostsByAuthorId field.
func (r *queryResolver) GetPostsByAuthorID(ctx context.Context, id string) (*model.Post, error) {
	authorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	post := model.Post{}

	if err := r.Database.Model(&model.Post{}).Preload("Author").Find(&post, "authorID = ?", uint(authorId)).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &post, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
