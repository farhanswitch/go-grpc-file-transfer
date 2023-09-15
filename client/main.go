package main

import (
	"context"
	"flag"
	"log"

	"github.com/farhanswitch/grpc-file/utilities/upload"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("Missing file path")
	}

	// Initialize GRPC Connection
	conn, err := grpc.Dial(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to GRPC Server.\nError: %s\n", err.Error())
	}
	defer conn.Close()

	// Start uploading file
	client := upload.NewClient(conn)
	name, err := client.Upload(context.Background(), flag.Arg(0))
	if err != nil {
		log.Fatalf("Error when uploading file.\nError: %s\n", err.Error())
	}
	log.Printf("File uploaded successfully.\nFile name: %s\n", name)
}
