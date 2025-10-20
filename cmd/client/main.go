package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-file-storage-go/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewFileServiceClient(conn)

	uploadedFilename := testUploadFile(client)
	time.Sleep(1 * time.Second)

	testListFiles(client)
	time.Sleep(1 * time.Second)

	testDownloadFile(client, uploadedFilename)
}

func testUploadFile(client proto.FileServiceClient) string {
	fmt.Println("=== Testing UploadFile ===")

	fileContent := []byte("Hello, this is a test file content!")

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("faled to create stream: %v", err)
	}

	err = stream.Send(&proto.UploadFileRequest{
		Data: &proto.UploadFileRequest_Info{
			Info: &proto.FileInfo{
				Filename:    "test.txt",
				ContentType: "text/plan",
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to send stream: %v", err)
	}

	chunkSize := 10
	for i := 0; i < len(fileContent); i += chunkSize {
		end := i + chunkSize
		if end > len(fileContent) {
			end = len(fileContent)
		}
		err = stream.Send(&proto.UploadFileRequest{
			Data: &proto.UploadFileRequest_ChunkData{
				ChunkData: fileContent[i:end],
			},
		})
		if err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}

	fmt.Printf("Upload successful: ID=%s, Filename=%s, Size=%d\n",
		response.GetId(), response.GetFilename(), response.GetSize())
	return response.GetFilename()
}

func testDownloadFile(client proto.FileServiceClient, filename string) {
	fmt.Printf("\n=== Testing DownloadFile for: %s===\n", filename)

	stream, err := client.DownloadFile(context.Background(), &proto.DownloadFileRequest{
		Filename: filename,
	})
	if err != nil {
		log.Printf("failed to download file: %v", err)
		return
	}

	var fileData []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed to receive chunk: %v", err)
			return
		}
		fileData = append(fileData, chunk.ChunkData...)
	}

	fmt.Printf("Downloaded %d bytes: %s\n", len(fileData), string(fileData))
}

func testListFiles(client proto.FileServiceClient) {
	fmt.Println("\n=== Testing ListFiles ===")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.ListFiles(ctx, &proto.ListFilesRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("failed to list files: %v", err)
		return
	}

	fmt.Printf("Found %d files\n", response.TotalCount)

	for _, file := range response.Files {
		fmt.Printf("- %s (size: %d bytes, created: %v)\n",
			file.Filename, file.Size, file.CreatedAt.AsTime())
	}
}
