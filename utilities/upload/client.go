package upload

import (
	"context"
	"io"
	"os"
	"time"

	filepb "github.com/farhanswitch/grpc-file/proto"
	"google.golang.org/grpc"
)

type Client struct {
	client filepb.FileServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: filepb.NewFileServiceClient(conn),
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

	// Maximum 1KB per stream
	buf := make([]byte, 1024)

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
