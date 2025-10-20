package grpc

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockDownloadFileStream struct {
	mock.Mock
	sentChunks []*proto.DownloadFileResponse
}

func (m *MockDownloadFileStream) Send(response *proto.DownloadFileResponse) error {
	m.sentChunks = append(m.sentChunks, response)
	args := m.Called(response)
	return args.Error(0)
}

func (m *MockDownloadFileStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *MockDownloadFileStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *MockDownloadFileStream) SetTrailer(metadata.MD) {
}

func (m *MockDownloadFileStream) Context() context.Context {
	return context.Background()
}

func (m *MockDownloadFileStream) SendMsg(interface{}) error {
	return nil
}

func (m *MockDownloadFileStream) RecvMsg(interface{}) error {
	return nil
}

func Test_DownloadFile_Success(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	testFileContent := "Hello, this is test file content for download!"
	testFile := &domain.File{
		ID:       "download-uuid",
		Filename: "test_download.txt",
		Size:     int64(len(testFileContent)),
		Path:     "/storage/test_download.txt",
	}

	fileReader := io.NopCloser(strings.NewReader(testFileContent))

	mockUseCase.On(
		"DownLoadFile",
		mock.Anything,
		"test_download.txt").
		Return(testFile, fileReader, nil)

	mockStream := new(MockDownloadFileStream)
	mockStream.On(
		"Send",
		mock.AnythingOfType("*proto.DownloadFileResponse")).
		Return(nil).
		Maybe()

	err := handler.DownloadFile(&proto.DownloadFileRequest{
		Filename: "test_download.txt",
	},
		mockStream)

	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)

	assert.Greater(t, len(mockStream.sentChunks), 0, "Should have sent all chunk")

	var allData []byte
	for _, chunk := range mockStream.sentChunks {
		allData = append(allData, chunk.ChunkData...)
	}
	assert.Equal(t, testFileContent, string(allData))
}

func Test_DownloadFile_FileNotFound(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	mockUseCase.On(
		"DownLoadFile",
		mock.Anything,
		"nonexistent.txt").
		Return(nil, nil, errors.New("file not found"))

	mockStream := new(MockDownloadFileStream)

	err := handler.DownloadFile(&proto.DownloadFileRequest{
		Filename: "nonexistent.txt",
	},
		mockStream)

	assert.Error(t, err)

	grpcStatus, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, grpcStatus.Code())
	assert.Equal(t, "file reader is nil", grpcStatus.Message())

	mockStream.AssertNotCalled(t, "Send", mock.AnythingOfType("*proto.DownloadFileResponse"))
	mockUseCase.AssertExpectations(t)
}

func Test_DownloadFile_EmptyFile(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	testFile := &domain.File{
		ID:       "empty-uuid",
		Filename: "empty.txt",
		Size:     0,
		Path:     "/storage/empty.txt",
	}
	fileReader := io.NopCloser(strings.NewReader(""))

	mockUseCase.On(
		"DownLoadFile",
		mock.Anything, "empty.txt").
		Return(testFile, fileReader, nil)

	mockStream := new(MockDownloadFileStream)

	err := handler.DownloadFile(&proto.DownloadFileRequest{
		Filename: "empty.txt",
	},
		mockStream)

	assert.NoError(t, err)
	mockUseCase.AssertExpectations(t)
	mockStream.AssertNotCalled(t, "Send")
}

func Test_DownloadFile_SendError(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	testFileContent := "Test content that will fail to send"
	testFile := &domain.File{
		ID:       "error-uuid",
		Filename: "error.txt",
		Size:     int64(len(testFileContent)),
		Path:     "/storage/error.txt",
	}

	fileReader := io.NopCloser(strings.NewReader(testFileContent))

	mockUseCase.On("DownLoadFile", mock.Anything, "error.txt").Return(testFile, fileReader, nil)

	mockStream := new(MockDownloadFileStream)
	mockStream.On(
		"Send",
		mock.AnythingOfType("*proto.DownloadFileResponse")).
		Return(errors.New("network error"))

	err := handler.DownloadFile(&proto.DownloadFileRequest{
		Filename: "error.txt",
	},
		mockStream)

	assert.Error(t, err)
	grpcStatus, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, grpcStatus.Code())

	mockUseCase.AssertExpectations(t)
	mockStream.AssertCalled(t, "Send", mock.AnythingOfType("*proto.DownloadFileResponse"))
}
