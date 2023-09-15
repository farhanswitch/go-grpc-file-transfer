package upload

import (
	"io"
	"math/rand"
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
