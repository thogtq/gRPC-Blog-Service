package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/thogtq/grpc-blog-service/m/v1/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct{}
type blogItem struct {
	ID       primitive.ObjectID `bjson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bjson:"content"`
	Title    string             `bson:"title"`
}

var collection *mongo.Collection

const secureConnection = false
const serverPort = "50051"
const serverAddress = "0.0.0.0:"

func main() {
	//logging line of code causes server crashing
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Mongodb connection setup
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("error when creating mongodb client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("can not connect to mongodb server")
	}
	collection = client.Database("mydb").Collection("blog")
	fmt.Println("Connected to MongoDB")

	listen, listenErr := net.Listen("tcp", serverAddress+serverPort)
	if listenErr != nil {
		log.Fatalf("error while listen tcp at %v%v", serverAddress, serverPort)
		return
	}

	//gRPC server setup
	serverOpts := []grpc.ServerOption{}
	if secureConnection {
		//establish TSL connection
	}
	serverControl := grpc.NewServer(serverOpts...)
	blogpb.RegisterBlogServiceServer(serverControl, &server{})
	fmt.Println("Blog Service started")
	//The *Server.Serve() function will block the program so we run it in a goroutine
	go func() {
		fmt.Println("Server started")
		serveErr := serverControl.Serve(listen)
		//Blocked if successfully serve
		if serveErr != nil {
			log.Fatalf("fail to serve server : %v", serveErr)
		}
	}()

	//Setup shutdown hook for server when Ctrl + C
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	//Block until shutdown signal is received
	<-shutdownCh
	//Prepare and shut down server
	fmt.Printf("\nStopping the server...\n")
	serverControl.Stop()
	fmt.Println("Closing the listener...")
	listen.Close()
	fmt.Println("Disconnecting from mongodb...")
	client.Disconnect(ctx)
	fmt.Println("Server shutdown reached")
}
