package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/anchalpuri/assignment/blog"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	// Create a new post
	newPost := &pb.Post{
		Title:   "My First Blog Post",
		Content: "This is the content of my first blog post",
		Author:  "John Doe",
		PubDate: fmt.Sprintf("%v", time.Now().Format("2006-01-02")),
	}
	createdPost, err := client.CreatePost(context.Background(), newPost)
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}
	fmt.Println("Created Post:", createdPost)

	// Get the created post
	postId := &pb.PostID{Id: createdPost.GetId()}
	retrievedPost, err := client.ReadPost(context.Background(), postId)
	if err != nil {
		log.Fatalf("failed to get post: %v", err)
	}
	fmt.Println("Retrieved Post:", retrievedPost)

	// Update the post (change title)
	retrievedPost.Title = "Updated Title"
	updatedPost, err := client.UpdatePost(context.Background(), retrievedPost)
	if err != nil {
		log.Fatalf("failed to update post: %v", err)
	}
	fmt.Println("Updated Post:", updatedPost)

	// Delete the post
	deletePostID := &pb.PostID{Id: retrievedPost.GetId()}
	response, err := client.DeletePost(context.Background(), deletePostID)
	if err != nil {
		log.Fatalf("failed to delete post: %v", err)
	}
	fmt.Println(response.Response)
}
