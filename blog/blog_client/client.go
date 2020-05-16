package main

import (
	"context"
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
	doCreateBlog(c)
}

func doCreateBlog(c blogpb.BlogServiceClient) {
	blog := &blogpb.CreateBlogReq{
		Blog: &blogpb.Blog{
			Author:  "pgandla",
			Title:   "Mongo post",
			Content: "Go & Mongo go perfectly fine",
		},
	}
	res, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	log.Printf("Blog posted: %v", res.GetBlog().Id)
}
