package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-lesson/pb"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("not connect %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("response form download: %v", res.GetData())
		log.Printf("response form download: %v", string(res.GetData()))
	}
}
