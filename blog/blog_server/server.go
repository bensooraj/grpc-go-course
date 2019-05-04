package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bensooraj/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var collection *mongo.Collection

// MongoDBCredentials ...
type MongoDBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type server struct {
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {

	// In case of crash, get the filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load MongoDB Credentials
	mongodbCredsFile, err := os.Open("mongodb-config/credentials.json")
	if err != nil {
		log.Fatalf("Unable to fetch mongoDB credentials")
	}
	var mongodbCredentials MongoDBCredentials
	decoder := json.NewDecoder(mongodbCredsFile)
	err = decoder.Decode(&mongodbCredentials)
	if err != nil {
		log.Fatalf("Unable to parse mongoDB credentials")
	}

	mongoDbConnectionURL := fmt.Sprintf("mongodb://%s:%s@ds149596.mlab.com:49596/mydb?authMechanism=SCRAM-SHA-1", mongodbCredentials.Username, mongodbCredentials.Password)
	fmt.Println("Connecting to MongoDB: ", mongoDbConnectionURL)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDbConnectionURL))
	if err != nil {
		log.Fatal(err)
	}
	mongoDBctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(mongoDBctx)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Server Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// Wait for Ctrl-C to exit
	osSignalChannel := make(chan os.Signal, 1)
	signal.Notify(osSignalChannel, os.Interrupt)

	// Block until a signal is received
	<-osSignalChannel

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(mongoDBctx)
	fmt.Println("End Of Program")

}
