package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bensooraj/grpc-go-course/blog/blogpb"
)

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot parse/convert to ObjectID: %v", err),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       objectID.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}
