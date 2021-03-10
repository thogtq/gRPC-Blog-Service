package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/thogtq/grpc-blog-service/m/v1/blogpb"
	"google.golang.org/grpc"
)

const secureConnection = false
const serverPort = "50051"
const serverAddress = "0.0.0.0:"

func main() {
	connOpts := []grpc.DialOption{}
	if secureConnection {
		//establish TSL connection
	} else {
		connOpts = append(connOpts, grpc.WithInsecure())
	}
	conn, dialErr := grpc.Dial(serverAddress+serverPort, connOpts...)
	if dialErr != nil {
		log.Fatalf("fail to connect server %v", dialErr)
	}
	defer conn.Close()
	service := blogpb.NewBlogServiceClient(conn)

	blogID := createBlog(service)
	readBlog(service, blogID)
	updateBlog(service, blogID)
	deleteBlog(service, blogID)
	listBlog(service)
}
func createBlog(service blogpb.BlogServiceClient) string {
	//create blog
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "QuocThong",
			Title:    "My first blog post",
			Content:  "Hello!. This is my first blog post",
		},
	}
	createBlogRes, createBlogErr := service.CreateBlog(context.Background(), req)
	if createBlogErr != nil {
		log.Printf("can not create blog : %v\n", createBlogErr)
	}
	fmt.Printf("Blog has been created : %v\n", createBlogRes.GetBlog())
	return createBlogRes.GetBlog().GetId()
}
func readBlog(service blogpb.BlogServiceClient, blogID string) {
	//read blog
	req := &blogpb.ReadBlogRequest{
		BlogId: blogID,
	}
	res, readErr := service.ReadBlog(context.Background(), req)
	if readErr != nil {
		log.Fatalf("can not read blog : %v\n", readErr)
	}
	fmt.Printf("Blog was read : %v\n", res.GetBlog())
}
func updateBlog(service blogpb.BlogServiceClient, blogID string) {
	//udpate blog
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       blogID,
			AuthorId: "Tran Quoc Thong",
			Title:    "My first blog updated",
			Content:  "Hello!",
		},
	}
	res, updateErr := service.UpdateBlog(context.Background(), req)
	if updateErr != nil {
		log.Fatalf("can not read blog : %v\n", updateErr)
	}
	fmt.Printf("Blog was updated : %v\n", res.GetBlog())
}
func deleteBlog(service blogpb.BlogServiceClient, blogID string) {
	//delete blog
	req := &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	}
	res, deleteErr := service.DeleteBlog(context.Background(), req)
	if deleteErr != nil {
		log.Fatalf("can not delete blog : %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted : %v\n", res.GetBlogId())
}
func listBlog(service blogpb.BlogServiceClient) {
	stream, err := service.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog : %v\n", err)
	}
	fmt.Printf("List of Blogs :\n")
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading blog : %v\n", err)
		}
		fmt.Printf("%v\n", res.GetBlog())
	}
}
