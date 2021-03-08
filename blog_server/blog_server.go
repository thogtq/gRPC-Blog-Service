package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/thogtq/grpc-blog-service/m/v1/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

const secureConnection = false
const serverPort = "50051"
const serverAddress = "0.0.0.0:"

func main() {
	//logging line of code causes server crashing
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	serverOpts := []grpc.ServerOption{}
	if secureConnection {
		//establish TSL connection
	}

	listen, listenErr := net.Listen("tcp", serverAddress+serverPort)
	if listenErr != nil {
		log.Fatalf("error while listen tcp at %v%v", serverAddress, serverPort)
		return
	}

	serverControl := grpc.NewServer(serverOpts...)
	blogpb.RegisterBlogServiceServer(serverControl, &server{})
	//The *Server.Serve() function will block the program so we run it in goroutine
	go func() {
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
	fmt.Println("Stopping the server...")
	serverControl.Stop()
	fmt.Println("Closing listener...")
	listen.Close()
	fmt.Println("Server shutdown reached")
}
