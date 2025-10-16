package main

import (
	"context"
	"flag"
	"log"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/ivagulin/grpc-flatbuffers-example/bookmarks"

	"google.golang.org/grpc"
)

type server struct{}

var flatAddr = flag.String("flatAddr", "0.0.0.0:50051", "gRPC server address")
var grpcAddr = flag.String("grpcAddr", "0.0.0.0:50052", "gRPC server address")
var cmd = flag.String("cmd", "last-added", "cmd")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*flatAddr, grpc.WithInsecure(), grpc.WithCodec(flatbuffers.FlatbuffersCodec{}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := bookmarks.NewBookmarksServiceClient(conn)

	if *cmd == "add" {

		if flag.NArg() < 2 {
			log.Fatalln("Insufficient args provided for add command..")
		}

		b := flatbuffers.NewBuilder(0)
		url := b.CreateString(flag.Arg(0))
		title := b.CreateString(flag.Arg(1))

		bookmarks.AddRequestStart(b)
		bookmarks.AddRequestAddUrl(b, url)
		bookmarks.AddRequestAddTitle(b, title)
		b.Finish(bookmarks.AddRequestEnd(b))

		_, err = client.Add(context.Background(), b)
		if err != nil {
			log.Fatalf("Retrieve flatClient failed: %v", err)
		}

	} else if *cmd == "last-added" {

		b := flatbuffers.NewBuilder(0)
		bookmarks.LastAddedRequestStart(b)
		b.Finish(bookmarks.LastAddedRequestEnd(b))

		out, err := client.LastAdded(context.Background(), b)
		if err != nil {
			log.Fatalf("Retrieve flatClient failed: %v", err)
		}

		log.Println("ID: ", string(out.Id()))
		log.Println("URL: ", string(out.Url()))
		log.Println("Title: ", string(out.Title()))

	}

	log.Println("SENT")

}
