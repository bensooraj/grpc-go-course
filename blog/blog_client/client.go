package main

import (
	"context"
	"fmt"
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
	// fmt.Println("Reading a blog")
	//
	blogData, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Fatalf("Couldn't find blog with the given blog ID")
	}
	fmt.Println("Blog Returned: ", blogData)
}
