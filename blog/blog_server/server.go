package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deepudoit/coolgo/gogrpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
}

var coll *mongo.Collection
var objectId primitive.ObjectID

type blogItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Author  string             `bson:"author_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}


func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()
	data := blogItem{
		Author:  blog.GetAuthor(),
		Title:   blog.GetTitle(),
		Content: blog.GetContent(),
	}
	res, err := coll.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error: %v", err))
	}
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Can't covert to ObjectID"))
	}
	return &blogpb.CreateBlogRes{
		Blog: &blogpb.Blog{
			Id:      objID.Hex(),
			Author:  blog.GetAuthor(),
			Title:   blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	fmt.Println("ReadBlog call...")
	blogID := req.GetBlogID()

	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot parse ID: %v", err))
	}
	data := &blogItem{}
	result := coll.FindOne(context.TODO(), bson.D{{"_id", oid}})
	if err := result.Decode(data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("No documents with the blog ID: %v", oid))
		}
	}
	res := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id: data.ID.Hex(),
			Author: data.Author,
			Title: data.Title,
			Content: data.Content,
		},
	}
	return res, nil
}

func main() {
	fmt.Println("Blog service...")
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//Mongo connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect mongo session: %v", err)
	}

	//Access DB objects
	coll = client.Database("mydb").Collection("blog")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	//Graceful shutdown of server and listener on Ctrl+C
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to server gRPC srvr: %v", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	s.Stop()
	err = client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("Failed to close Mongo connection: %v", err)
	}
	err = lis.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Stopping server and listener...")
}
