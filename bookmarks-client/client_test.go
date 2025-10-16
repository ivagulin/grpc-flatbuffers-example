package main

import (
	"context"
	"log"
	"sync"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/ivagulin/grpc-flatbuffers-example/bookmarks"
	"github.com/ivagulin/grpc-flatbuffers-example/bookmarkspb"
	"google.golang.org/grpc"
)

var flatClient = sync.OnceValue(func() bookmarks.BookmarksServiceClient {
	conn, err := grpc.Dial(*flatAddr, grpc.WithInsecure(), grpc.WithCodec(flatbuffers.FlatbuffersCodec{}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//defer conn.Close()
	client := bookmarks.NewBookmarksServiceClient(conn)
	return client
})

var onceBuilder = sync.OnceValue(func() *flatbuffers.Builder {
	return flatbuffers.NewBuilder(0)
})

var grpcClient = sync.OnceValue(func() bookmarkspb.BookmarksServiceClient {
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//defer conn.Close()
	return bookmarkspb.NewBookmarksServiceClient(conn)
})

func BenchmarkFlat(b *testing.B) {
	fb := onceBuilder()
	fc := flatClient()

	for i := 0; i < b.N; i++ {
		bookmarks.LastAddedRequestStart(fb)
		fb.Finish(bookmarks.LastAddedRequestEnd(fb))

		_, err := fc.LastAdded(context.Background(), fb)
		if err != nil {
			log.Fatalf("Retrieve flatClient failed: %v", err)
		}
	}
}

func BenchmarkGRPC(b *testing.B) {
	gc := grpcClient()
	for i := 0; i < b.N; i++ {
		_, err := gc.LastAdded(context.Background(), &bookmarkspb.LastAddedRequest{})
		if err != nil {
			log.Fatalf("Retrieve grpcClient failed: %v", err)
		}
	}
}
