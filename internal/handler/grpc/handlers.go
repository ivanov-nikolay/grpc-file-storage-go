package grpc

import (
	"context"
	"io"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *fileHandler) UploadFile(stream proto.FileService_UploadFileServer) error {
	var fileInfo *proto.FileInfo
	var data []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		switch x := req.Data.(type) {
		case *proto.UploadFileRequest_Info:
			fileInfo = x.Info
		case *proto.UploadFileRequest_ChunkData:
			data = append(data, x.ChunkData...)
		}
	}

	if fileInfo == nil {
		return status.Error(codes.InvalidArgument, "file info is required")
	}
	file, err := h.fileUseCase.UploadFile(
		stream.Context(),
		fileInfo.Filename,
		&bytesReader{data: data},
	)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&proto.UploadFileResponse{
		Id:       file.ID,
		Filename: file.Filename,
		Size:     uint32(file.Size),
	})
}

func (h *fileHandler) DownloadFile(*proto.DownloadFileRequest, proto.FileService_DownloadFileServer) error {
	return nil
}

func (h *fileHandler) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	return &proto.ListFilesResponse{}, nil
}

type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	n = copy(p, r.data[r.pos:])
	r.pos += n

	return n, nil
}
