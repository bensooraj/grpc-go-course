package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/bensooraj/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial/connect: %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create Blog
	fmt.Println("Creating a blog")
	blog := &blogpb.Blog{
		AuthorId: "Sooraj",
		Title:    "My Third Blog",
		Content:  "Third Blog Contents",
	}
	createdBlog, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected Error: %v", err)
	}
	fmt.Printf("Blog created: %v\n\n", createdBlog)
	blogID := createdBlog.GetBlog().GetId()

	// Read Blog
	fmt.Println("Reading a blog")

	blogData, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Fatalf("Couldn't find blog with the given blog ID")
	}
	fmt.Println("Blog Returned: ", blogData)

	// Update Blog
	fmt.Println("Updating a blog")
	updatedBlog, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5ccdf25ab7cd8d281541b83f",
			AuthorId: "Ben",
			Title:    "My First Blog",
			Content:  "First Blog Contents - 1",
		},
	})
	if err != nil {
		log.Fatalf("Couldn't find blog with the given blog ID")
	}
	fmt.Println("Blog Updated and Returned: ", updatedBlog)

	// Delete Blog
	fmt.Println("Updating a blog")
	deletedBlogID, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "5ccec7abd9cf356dc5dca4e0",
	})
	if err != nil {
		log.Fatalf("Couldn't find blog with the given blog ID")
	}
	fmt.Println("Blog Deleted: ", deletedBlogID)

	// List Blog
	fmt.Println("Updating a blog")
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog streaming RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something went wrong while iterating through the stream: %v", err)
		}

		fmt.Println(res.GetBlog())
	}
}
