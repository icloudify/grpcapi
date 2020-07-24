package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {
}

func (s server) ListBlog(request *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List Blog called ...")

	res, err := collection.Find(context.Background(), primitive.D{{}})
	defer res.Close(context.Background())
	if err != nil {
		return status.Errorf(codes.Internal, "Internal Error %v", err)
	}
	for res.Next(context.Background()) {
		data := &blogItem{}
		res.Decode(data)
		stream.Send(&blogpb.ListBlogResponse{
			Blog: dataToBlogPb(data),
		})
	}

	return nil
}

func (s server) DeleteBlog(ctx context.Context, request *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete Blog called ...")
	oid, err := primitive.ObjectIDFromHex(request.GetBlogId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument Error %v", err)
	}

	filter := bson.M{"_id": oid}

	_, error := collection.DeleteOne(context.Background(), filter)

	if error != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error %v", error)
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: request.GetBlogId(),
	}, nil
}

func (s server) UpdateBlog(ctx context.Context, request *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("CreateBlog called ...")
	blog := request.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument Error %v", err)
	}
	data := &blogItem{}
	data.ID = oid
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()
	data.AuthorID = blog.GetAuthorId()

	filter := bson.M{"_id": oid}

	_, err = collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error %v", err)
	}

	fmt.Println("Object id ", oid)
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (s server) ReadBlog(ctx context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read Blog called ...")
	oid, err := primitive.ObjectIDFromHex(request.GetBlogId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument Error %v", err)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)

	if res.Err() != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error %v", res.Err())
	}
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func (s server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("CreateBlog called ...")
	blog := request.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID %v", err)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

var collection *mongo.Collection

func main() {
	// in case of crash
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to mongodb")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error while creating mongodb client ", err)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal("Error while connecting mongodb ", err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("***********************\n  GRPC Server Started  \n***********************\n")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Error while starting srver", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("Faileed to serve", err)
		}
	}()

	// wait for ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	fmt.Println("Stopping the serve")
	s.Stop()
	lis.Close()
	client.Disconnect(context.TODO())
	fmt.Println("Closing mongo db connection")
	fmt.Println("closed the program")

}
