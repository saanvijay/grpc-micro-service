package main

import (
	"fmt"
	"log"
	"scmpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	port := 50051
	tls := true
	opts := grpc.WithInsecure()

	// With TLS/ without TLS ???
	if tls {
		certFile := "../TLS/HMSca.crt" // CA trust certificate
		creds, tlsErr := credentials.NewClientTLSFromFile(certFile, "")
		if tlsErr != nil {
			log.Fatalf("Error loading CA trust certificate : %v", tlsErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	connection, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts)
	if err != nil {
		log.Fatalf("Unable to connect localhost:%d %v", port, err)
	}
	defer connection.Close()

	client := scmpb.NewScmServiceClient(connection)

	TestNetwork(client)
	//TestSupplier(client)
	//TestManufacturer(client)
	//TestTransporter(client)
	//TestAgent(client)

}
