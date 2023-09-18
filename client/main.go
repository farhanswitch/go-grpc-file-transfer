package main

import (
	"context"
	"flag"
	"log"

	"github.com/farhanswitch/grpc-file/utilities/storage"
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

	var destinationFolder string = "./data/"
	// Start uploading file
	client := upload.NewClient(conn, storage.NewStorage(destinationFolder))
	name, err := client.Upload(context.Background(), flag.Arg(0))
	if err != nil {
		log.Fatalf("Error when uploading file.\nError: %s\n", err.Error())
	}
	log.Printf("File uploaded successfully.\nFile name: %s\n", name)

	var fileNameToDownload string = "vuoixIBaDp.png"

	err = client.Write(context.Background(), fileNameToDownload)
	if err != nil {
		log.Printf("Failed to download file.\nError: %s\n", err.Error())
	}
	log.Printf("Success download file %s.\nSaved at %s", fileNameToDownload, destinationFolder+fileNameToDownload)

}
