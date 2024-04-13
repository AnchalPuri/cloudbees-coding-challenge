package main

import (
	"context"
	"testing"

	pb "github.com/anchalpuri/assignment/blog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testPost = "Test Post"
const content = "This is a test post content"

func TestCreatePost(t *testing.T) {
	server := &BlogServer{}
	newPost := &pb.Post{
		Title:   testPost,
		Content: content,
		Author:  "Tester",
	}

	createdPost, err := server.CreatePost(context.Background(), newPost)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error creating post: %v", err)
	}
	if createdPost.GetId() == "" {
		t.Errorf("Expected generated ID, got empty ID")
	}
	if len(posts) != 1 {
		t.Errorf("Expected one post in storage, got %d", len(posts))
	}
}

func TestCreatePostInvalidRequest(t *testing.T) {
	server := &BlogServer{}
	invalidPost := &pb.Post{
		Content: content,
		Author:  "Tester",
	}

	// Act
	_, err := server.CreatePost(context.Background(), invalidPost)

	// Assert
	if err == nil {
		t.Errorf("Expected error for invalid post")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error code, got %v", status.Code(err))
	}
}

func TestReadPost(t *testing.T) {
	server := &BlogServer{}
	newPost := &pb.Post{
		Title:   testPost,
		Content: content,
		Author:  "Tester",
	}
	createdPost, err := server.CreatePost(context.Background(), newPost)
	if err != nil {
		t.Fatal(err)
	}

	retrievedPost, err := server.ReadPost(context.Background(), &pb.PostID{Id: createdPost.GetId()})

	if err != nil {
		t.Errorf("Unexpected error retrieving post: %v", err)
	}
	if retrievedPost.GetId() != createdPost.GetId() {
		t.Errorf("Retrieved post ID mismatch")
	}
}

func TestReadPostNotFound(t *testing.T) {

	server := &BlogServer{}
	invalidID := "non-existent-id"

	_, err := server.ReadPost(context.Background(), &pb.PostID{Id: invalidID})

	if err == nil {
		t.Errorf("Expected error for non-existent post ID")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error code, got %v", status.Code(err))
	}
}

func TestUpdatePost(t *testing.T) {
	server := &BlogServer{}
	newPost := &pb.Post{
		Title:   testPost,
		Content: content,
		Author:  "Tester",
	}
	createdPost, err := server.CreatePost(context.Background(), newPost)
	if err != nil {
		t.Fatal(err)
	}

	createdPost.Title = "Updated Title"

	updatedPost, err := server.UpdatePost(context.Background(), createdPost)

	if err != nil {
		t.Errorf("Unexpected error updating post: %v", err)
	}
	if updatedPost.GetTitle() != createdPost.GetTitle() { // Check if title is updated
		t.Errorf("Expected updated title in response, got %s", updatedPost.GetTitle())
	}
}
