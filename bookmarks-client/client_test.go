package main

import (
	"context"
	"log"
	"sync"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/ivagulin/grpc-flatbuffers-example/bookmarks"
	"google.golang.org/grpc"
)

var client = sync.OnceValue(func() bookmarks.BookmarksServiceClient {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithCodec(flatbuffers.FlatbuffersCodec{}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//defer conn.Close()
	client := bookmarks.NewBookmarksServiceClient(conn)
	return client
})

func BenchmarkLastAdded(b *testing.B) {
	fb := flatbuffers.NewBuilder(0)

	for i := 0; i < b.N; i++ {
		bookmarks.LastAddedRequestStart(fb)
		fb.Finish(bookmarks.LastAddedRequestEnd(fb))

		_, err := client().LastAdded(context.Background(), fb)
		if err != nil {
			log.Fatalf("Retrieve client failed: %v", err)
		}
	}
}
