package upload

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	filepb "github.com/farhanswitch/grpc-file/proto"
	"github.com/farhanswitch/grpc-file/utilities/storage"
	"google.golang.org/grpc"
)

type Client struct {
	storage storage.Storage
	client  filepb.FileServiceClient
}

func NewClient(conn grpc.ClientConnInterface, storage storage.Storage) Client {
	return Client{
		client:  filepb.NewFileServiceClient(conn),
		storage: storage,
	}
}
func (c Client) Write(con context.Context, fileName string) error {
	var file *storage.File = storage.NewFile(fileName)
	ctx, cancel := context.WithDeadline(con, time.Now().Add(20*time.Second))
	defer cancel()
	stream, err := c.client.Download(ctx, &filepb.DownloadRequest{
		Name: fileName,
	})
	if err != nil {
		return err
	}
	for {
		req, err := stream.Recv()
		if file == nil {
			log.Println("Empty")
			file = storage.NewFile(req.GetName())
		}
		if err == io.EOF {
			if err := c.storage.Store(file); err != nil {
				log.Printf("Error when save downloaded file.\nError: %s\n", err.Error())
				return err
			}
			return stream.CloseSend()
		}
		if err != nil {
			return err
		}
		if err := file.Write(req.GetChunk()); err != nil {
			return err
		}
	}

}
func (c Client) Upload(ctx context.Context, file string) (string, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(20*time.Second))
	defer cancel()

	stream, err := c.client.Upload(ctx)
	if err != nil {
		return "", err
	}
	fil, err := os.Open(file)
	if err != nil {
		return "", err
	}

	// Maximum 10KB per stream
	buf := make([]byte, 10*1024)

	for {
		num, err := fil.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if err := stream.Send(&filepb.UploadRequest{Chunk: buf[:num]}); err != nil {
			return "", err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}
	return res.GetName(), nil
}
