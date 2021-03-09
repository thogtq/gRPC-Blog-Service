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

	//createBlog(service)
	//readBlog(service)
	updateBlog(service)

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
		log.Printf("can not create blog : %v\n", createBlogErr)
	}
	fmt.Printf("Blog has been created : %v\n", createBlogRes.GetBlog())
}
func readBlog(service blogpb.BlogServiceClient) {
	//read blog
	req := &blogpb.ReadBlogRequest{
		BlogId: "604750491ab91f87c478f3db",
	}
	res, readErr := service.ReadBlog(context.Background(), req)
	if readErr != nil {
		log.Fatalf("can not read blog : %v\n", readErr)
	}
	fmt.Printf("Blog was read : %v\n", res.GetBlog())
}
func updateBlog(service blogpb.BlogServiceClient) {
	//udpate blog
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "604750491ab91f87c478f3db",
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
