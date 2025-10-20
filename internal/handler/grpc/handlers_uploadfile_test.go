package grpc

import (
	"context"
	"io"
	"testing"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockFileUseCase struct {
	mock.Mock
}

func (m *MockFileUseCase) UploadFile(ctx context.Context, filename string, data io.Reader) (*domain.File, error) {
	args := m.Called(ctx, filename, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *MockFileUseCase) DownLoadFile(ctx context.Context, filename string) (*domain.File, io.Reader, error) {
	args := m.Called(ctx, filename)
	if args.Get(0) == nil {
		return nil, nil, args.Error(1)
	}

	return args.Get(0).(*domain.File), args.Get(1).(io.Reader), args.Error(2)
}

func (m *MockFileUseCase) ListFiles(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	args := m.Called(ctx, page, pageSize)

	return args.Get(0).(*domain.FileList), args.Error(1)
}

type MockUploadFileStream struct {
	mock.Mock
	requests []*proto.UploadFileRequest
	response *proto.UploadFileResponse
}

func (m *MockUploadFileStream) SendAndClose(response *proto.UploadFileResponse) error {
	m.response = response
	args := m.Called(response)

	return args.Error(0)
}

func (m *MockUploadFileStream) Recv() (*proto.UploadFileRequest, error) {
	if len(m.requests) == 0 {
		return nil, io.EOF
	}
	req := m.requests[0]
	m.requests = m.requests[1:]

	return req, nil
}

func (m *MockUploadFileStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *MockUploadFileStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *MockUploadFileStream) SetTrailer(metadata.MD) {
}

func (m *MockUploadFileStream) Context() context.Context {
	return context.Background()
}

func (m *MockUploadFileStream) SendMsg(interface{}) error {
	return nil
}

func (m *MockUploadFileStream) RecvMsg(interface{}) error {
	return nil
}

func Test_UploadFile_Success(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	testFileContent := "Hello, this is test file content!"
	expectedFile := &domain.File{
		ID:       "test-uuid",
		Filename: "test_123.txt",
		Size:     int64(len(testFileContent)),
		Path:     "/storage/test_123.txt",
	}

	mockUseCase.On(
		"UploadFile",
		mock.Anything,
		"test.txt",
		mock.AnythingOfType("*grpc.bytesReader")).
		Return(expectedFile, nil).
		Run(func(args mock.Arguments) {
			reader := args.Get(2).(io.Reader)
			data, _ := io.ReadAll(reader)
			assert.Equal(t, testFileContent, string(data))
		})

	mockStream := new(MockUploadFileStream)
	mockStream.requests = []*proto.UploadFileRequest{
		{
			Data: &proto.UploadFileRequest_Info{
				Info: &proto.FileInfo{
					Filename:    "test.txt",
					ContentType: "text/plain",
				},
			},
		},
		{
			Data: &proto.UploadFileRequest_ChunkData{
				ChunkData: []byte(testFileContent),
			},
		},
	}
	mockStream.On(
		"SendAndClose",
		mock.AnythingOfType("*proto.UploadFileResponse")).
		Return(nil)

	err := handler.UploadFile(mockStream)
	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
	mockStream.AssertExpectations(t)

	if mockStream.response != nil {
		assert.Equal(t, "test-uuid", mockStream.response.Id)
		assert.Equal(t, "test_123.txt", mockStream.response.Filename)
		assert.Equal(t, uint32(len(testFileContent)), mockStream.response.Size)
	}
}

func Test_UploadFile_NoFileInfo(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	mockStream := new(MockUploadFileStream)
	mockStream.requests = []*proto.UploadFileRequest{
		{
			Data: &proto.UploadFileRequest_ChunkData{
				ChunkData: []byte("test data"),
			},
		},
	}

	err := handler.UploadFile(mockStream)

	assert.Error(t, err)

	grpcStatus, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
	assert.Equal(t, "file info is required", grpcStatus.Message())

	mockStream.AssertNotCalled(t, "SendAndClose", mock.AnythingOfType("*proto.UploadFileResponse"))

	mockUseCase.AssertNotCalled(t, "UploadFile")
}

func Test_UploadFile_UseCaseError(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	mockUseCase.On(
		"UploadFile",
		mock.Anything,
		"test.txt",
		mock.AnythingOfType("*grpc.bytesReader")).
		Return(nil, assert.AnError)

	mockStream := new(MockUploadFileStream)
	mockStream.requests = []*proto.UploadFileRequest{
		{
			Data: &proto.UploadFileRequest_Info{
				Info: &proto.FileInfo{
					Filename:    "test.txt",
					ContentType: "text/plain",
				},
			},
		},
		{
			Data: &proto.UploadFileRequest_ChunkData{
				ChunkData: []byte("test data"),
			},
		},
	}

	err := handler.UploadFile(mockStream)
	assert.Error(t, err)

	grpcStatus, ok := status.FromError(err)
	assert.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.Internal, grpcStatus.Code())
	assert.Contains(t, grpcStatus.Message(), "assert.AnError general error for testing")

	mockUseCase.AssertExpectations(t)
	mockStream.AssertNotCalled(t, "SendAndClose", mock.AnythingOfType("*proto.UploadFileResponse"))
}
