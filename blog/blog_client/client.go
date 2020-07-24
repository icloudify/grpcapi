package main

import (
	"context"
	"fmt"
	"github.com/ravindra031/grpcapi/blog/blogpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}

	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)
	//doUnary(c)
	//doUnaryUpdate(c)
	//doUnaryDelete(c)
	doUnaryList(c)
}

func doUnary(c blogpb.BlogServiceClient) {
	req := blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Ravindra",
			Title:    "Rock it",
			Content:  "Rock content",
		},
	}
	res, err := c.CreateBlog(context.Background(), &req)

	if err != nil {
		log.Fatal("Could not connect to grpc!", err)
	}
	fmt.Println(res)

	fmt.Println("Reading the blog")

	data := &blogpb.ReadBlogRequest{
		BlogId: "5f1975903d5d2bce85850a82",
	}

	rbres, err := c.ReadBlog(context.Background(), data)

	if err != nil {
		log.Fatal("Could not read data from blog!", err)
	}

	fmt.Printf("Read the blog :\n ID %v\n Title %v\n Content %v\n Author Id %v\n", rbres.GetBlog().GetId(), rbres.GetBlog().GetTitle(), rbres.GetBlog().GetContent(), rbres.GetBlog().GetAuthorId())

}

func doUnaryUpdate(c blogpb.BlogServiceClient) {
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5f1975903d5d2bce85850a82",
			AuthorId: "Rock Ravindra",
			Title:    "Rock it",
			Content:  "Rock content update",
		},
	}

	rbres, err := c.UpdateBlog(context.Background(), req)

	if err != nil {
		log.Fatal("Could not read data from blog!", err)
	}

	fmt.Printf("Read the blog :\n ID %v\n Title %v\n Content %v\n Author Id %v\n", rbres.GetBlog().GetId(), rbres.GetBlog().GetTitle(), rbres.GetBlog().GetContent(), rbres.GetBlog().GetAuthorId())

}

func doUnaryDelete(c blogpb.BlogServiceClient) {
	req := &blogpb.DeleteBlogRequest{
		BlogId: "5f1975903d5d2bce85850a82",
	}

	rbres, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		log.Fatal("Could not read data from blog!", err)
	}

	fmt.Printf("Delete the blog :ID %v\n", rbres.GetBlogId())
}

func doUnaryList(c blogpb.BlogServiceClient) {
	// list Blogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
