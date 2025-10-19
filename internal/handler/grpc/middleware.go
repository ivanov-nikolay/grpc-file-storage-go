package grpc

import (
	"context"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConcurrencyLimiter struct {
	uploadSem   *semaphore.Weighted
	downloadSem *semaphore.Weighted
	listSem     *semaphore.Weighted
}

func NewConcurrencyLimiter(uploadLimit, downloadLimit, listLimit int64) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		uploadSem:   semaphore.NewWeighted(uploadLimit),
		downloadSem: semaphore.NewWeighted(downloadLimit),
		listSem:     semaphore.NewWeighted(listLimit),
	}
}

func (l *ConcurrencyLimiter) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/file_service.FileService/ListFiles" {
			if !l.listSem.TryAcquire(1) {
				return nil, status.Error(codes.ResourceExhausted,
					"too many concurrent list requests (max 100)")
			}
			defer l.listSem.Release(1)
		}

		return handler(ctx, req)
	}
}

func (l *ConcurrencyLimiter) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var sem *semaphore.Weighted

		switch info.FullMethod {
		case "/file_service.FileService/UploadFile":
			sem = l.uploadSem
		case "/file_service.FileService/DownloadFile":

			sem = l.downloadSem
		default:
			return handler(srv, stream)
		}

		if !sem.TryAcquire(1) {
			return status.Error(codes.ResourceExhausted,
				"too many concurrent requests (max 10)")
		}
		defer sem.Release(1)

		return handler(srv, stream)
	}
}
