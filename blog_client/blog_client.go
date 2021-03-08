package main

import (
	"context"
	"fmt"
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
	
	createBlog(service)

}
func createBlog(service blogpb.BlogServiceClient) {
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
		log.Fatalf("can not create blog : %v", createBlogErr)
	}
	fmt.Printf("Blog has been created : %v", createBlogRes.GetBlog())
}
