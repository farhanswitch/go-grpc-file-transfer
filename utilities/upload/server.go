package upload

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/farhanswitch/grpc-file/utilities/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	filepb "github.com/farhanswitch/grpc-file/proto"
)

type Server struct {
	storage storage.Storage
	filepb.UnimplementedFileServiceServer
}

func RandomString(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var charset string = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := make([]byte, length)
	for i := range str {
		str[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(str)
}
func NewServer(storage storage.Storage) Server {
	return Server{storage, filepb.UnimplementedFileServiceServer{}}
}
func (s Server) Download(req *filepb.DownloadRequest, stream filepb.FileService_DownloadServer) error {
	name := s.storage.Dir + req.GetName()
	fil, err := os.Open(name)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}
	defer fil.Close()
	// Maximum 10KB per stream
	buf := make([]byte, 10*1024)

	for {
		num, err := fil.Read(buf)
		// log.Println(buf[:num])

		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&filepb.DownloadResponse{
			Name:  name,
			Chunk: buf[:num],
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
func (s Server) Upload(stream filepb.FileService_UploadServer) error {
	name := RandomString(10) + ".png"
	file := storage.NewFile(name)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := s.storage.Store(file); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			return stream.SendAndClose(&filepb.UploadResponse{Name: name})
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		if err := file.Write(req.GetChunk()); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
