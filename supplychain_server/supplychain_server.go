package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"supplychainpb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type server struct{}

var mClient *mongo.Client

func main() {

	grpcServerPort := 50051
	grpcServerHost := "localhost"

	mongodbServerPort := 27017
	mongodbServerHost := "localhost"

	tls := true
	opts := []grpc.ServerOption{}
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// MongoDB Client
	mClient, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", mongodbServerHost, mongodbServerPort)))
	if err != nil {
		log.Fatal(err)
	}
	// Connect Client to db server
	err = mClient.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server Listening on Port %d ...\n", grpcServerPort)

	// Lets focus on Listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcServerHost, grpcServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	// With TLS/ without TLS ???
	if tls {
		certFile := "../TLS/scmserver.crt"
		keyFile := "../TLS/scmserver.pem" // .key (private key) file will not understandable by GRPC (.pem can be)
		creds, tlsErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if tlsErr != nil {
			log.Fatalf("Failed loading server certificates %v\n", tlsErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// Create GRPC Server
	grpcServer := grpc.NewServer(opts...)
	supplychainpb.RegisterScmServiceServer(grpcServer, &server{})

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Bind the port to GRPC Server
	err1 := grpcServer.Serve(listener)
	if err1 != nil {
		log.Fatalf("Failed to server %v\n", err1)
	}

}
