package main

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("List Files")

	dir := "/Users/abekouhei/Documents/abeshi-projects/grpc-lesson/storage"

	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}
	return res, nil
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was")

	filename := req.GetFilename()
	path := "/Users/abekouhei/Documents/abeshi-projects/grpc-lesson/storage" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	buf := make([]byte, 5)

	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err != io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("reveived data(bytes): %v", data)
		log.Printf("reveived data(bytes): %v", string(data))
		buf.Write(data)
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is run")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Faild to serve: %v", err)
	}
}
