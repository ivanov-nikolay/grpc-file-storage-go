package grpc

import (
	"context"
	"testing"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ListFiles(t *testing.T) {
	mockUseCase := new(MockFileUseCase)
	handler := NewFileHandler(mockUseCase)

	expectedFiles := &domain.FileList{
		Files: []domain.File{
			{
				ID:       "1",
				Filename: "file1.txt",
				Size:     100,
			},
		},
		Total: 1,
	}

	mockUseCase.On("ListFiles", mock.Anything, 1, 10).Return(expectedFiles, nil)

	resp, err := handler.ListFiles(context.Background(), &proto.ListFilesRequest{
		Page:     1,
		PageSize: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Files))
	assert.Equal(t, "file1.txt", resp.Files[0].Filename)
	mockUseCase.AssertExpectations(t)
}
