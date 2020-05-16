package main

import (
	"context"
	"fmt"
	"log"

	"github.com/deepudoit/coolgo/gogrpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect gRPC: %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	blogID := doCreateBlog(c)
	req := &blogpb.ReadBlogReq{
		BlogID: blogID,
	}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error reading blog posts: %s\n", err)
	}
	fmt.Printf("Blog : %v", res.GetBlog())
}

func doCreateBlog(c blogpb.BlogServiceClient) string {
	blog := &blogpb.CreateBlogReq{
		Blog: &blogpb.Blog{
			Author:  "pgandla",
			Title:   "gRPC Go Micro",
			Content: "Perfect way to communicate between services",
		},
	}
	res, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	log.Printf("Blog posted: %v", res.GetBlog().Id)
	return res.GetBlog().Id
}
