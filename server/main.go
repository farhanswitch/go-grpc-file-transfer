package main

import (
	"log"
	"net"

	"github.com/farhanswitch/grpc-file/utilities/storage"
	"github.com/farhanswitch/grpc-file/utilities/upload"

	"google.golang.org/grpc"

	filepb "github.com/farhanswitch/grpc-file/proto"
)

func main() {
	// Initialize GRPC Server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Cannot listening to port 50052.\nError: %s\n", err.Error())
	}
	defer lis.Close()

	fileUpload := upload.NewServer(storage.NewStorage("./tmp/"))

	rpcServer := grpc.NewServer()

	// Register and start RPC server
	filepb.RegisterFileServiceServer(rpcServer, fileUpload)

	log.Printf("GRPC Server listening on %v\n", lis.Addr().String())
	err = rpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Cannot starting GRPC Server.\nError: %s\n", err.Error())
	}

}
