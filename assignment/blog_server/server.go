package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/anchalpuri/assignment/blog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const postNotFoundErrMsg = "Post not found: %s"

// In-memory storage for posts
var posts map[string]*pb.Post = make(map[string]*pb.Post)

// BlogServer implements the BlogService server interface
type BlogServer struct {
	pb.UnimplementedBlogServiceServer
}

// CreatePost implements the RPC method to create a new post
func (s *BlogServer) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	// Generate a unique ID
	newID := fmt.Sprintf("post-%d", len(posts)+1)
	req.Id = newID

	// Validate and store the post
	if err := validatePost(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Post: %v", err)
	}
	posts[newID] = req

	// Return the created post
	return req, nil
}

// GetPost retrieves a post by its ID
func (s *BlogServer) ReadPost(ctx context.Context, req *pb.PostID) (*pb.Post, error) {
	postID := req.GetId()
	post, found := posts[postID]
	if !found {
		return nil, status.Errorf(codes.NotFound, postNotFoundErrMsg, postID)
	}
	return post, nil
}

// UpdatePost updates an existing post
func (s *BlogServer) UpdatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	postID := req.GetId()
	post, found := posts[postID]
	if !found {
		return nil, status.Errorf(codes.NotFound, postNotFoundErrMsg, postID)
	}

	// Update the post in the map
	post.Title = req.Title
	post.Content = req.Content
	post.Author = req.Author
	post.Tags = req.Tags

	posts[postID] = post

	return post, nil
}

// DeletePost deletes a post by its ID
func (s *BlogServer) DeletePost(ctx context.Context, req *pb.PostID) (*pb.Response, error) {
	postID := req.GetId()
	_, found := posts[postID]
	if !found {
		return &pb.Response{Response: "Post Not Found!"}, status.Errorf(codes.NotFound, postNotFoundErrMsg, postID)
	}

	delete(posts, postID)

	return &pb.Response{Response: "Post Successfully Deleted!"}, nil // Empty message for successful deletion
}

// validatePost performs basic validation on the Post object
func validatePost(post *pb.Post) error {
	if post.Title == "" || post.Content == "" || post.Author == "" {
		return fmt.Errorf("title, content, and author are required fields")
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterBlogServiceServer(server, &BlogServer{})

	log.Println("Server listening on port 8080")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
