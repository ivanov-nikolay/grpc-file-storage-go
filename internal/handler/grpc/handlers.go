package grpc

import (
	"context"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/usecase"
)

type fileHandler struct {
	proto.UnimplementedFileServiceServer
	fileUseCase usecase.FileUseCase
}

func NewFileHandler(fileUseCase usecase.FileUseCase) proto.FileServiceServer {
	return &fileHandler{
		fileUseCase: fileUseCase,
	}
}

func (h *fileHandler) UploadFile(proto.FileService_UploadFileServer) error {
	return nil
}

func (h *fileHandler) DownloadFile(*proto.DownloadFileRequest, proto.FileService_DownloadFileServer) error {
	return nil
}

func (h *fileHandler) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	return &proto.ListFilesResponse{}, nil
}
