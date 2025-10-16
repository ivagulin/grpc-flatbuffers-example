package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/ivagulin/grpc-flatbuffers-example/bookmarkspb"
	context "golang.org/x/net/context"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/ivagulin/grpc-flatbuffers-example/bookmarks"

	"google.golang.org/grpc"
)

type FlatServer struct {
	bookmarks.UnimplementedBookmarksServiceServer
	id        int
	lastTitle string
	lastURL   string
}

var flatAddr = "0.0.0.0:50051"
var grpcAddr = "0.0.0.0:50052"

func (s *FlatServer) Add(context context.Context, in *bookmarks.AddRequest) (*flatbuffers.Builder, error) {
	//log.Println("Add called...")

	s.id++
	s.lastTitle = string(in.Title())
	s.lastURL = string(in.Url())

	b := flatbuffers.NewBuilder(0)
	bookmarks.AddResponseStart(b)
	b.Finish(bookmarks.AddResponseEnd(b))
	return b, nil
}

func (s *FlatServer) LastAdded(context context.Context, in *bookmarks.LastAddedRequest) (*flatbuffers.Builder, error) {
	//log.Println("LastAdded called...")

	b := flatbuffers.NewBuilder(0)
	id := b.CreateString(strconv.Itoa(s.id))
	title := b.CreateString(s.lastTitle)
	url := b.CreateString(s.lastURL)

	bookmarks.LastAddedResponseStart(b)
	bookmarks.LastAddedResponseAddId(b, id)
	bookmarks.LastAddedResponseAddTitle(b, title)
	bookmarks.LastAddedResponseAddUrl(b, url)
	b.Finish(bookmarks.LastAddedResponseEnd(b))
	return b, nil
}

type GRPCServer struct {
	bookmarkspb.UnimplementedBookmarksServiceServer
	id        int64
	lastTitle string
	lastURL   string
}

func (s *GRPCServer) Add(ctx context.Context, req *bookmarkspb.AddRequest) (*bookmarkspb.AddResponse, error) {
	s.id++
	s.lastTitle = req.GetTitle()
	s.lastURL = req.GetURL()
	return &bookmarkspb.AddResponse{}, nil
}

func (s *GRPCServer) LastAdded(context.Context, *bookmarkspb.LastAddedRequest) (*bookmarkspb.LastAddedResponse, error) {
	//log.Println("LastAdded called...")
	return &bookmarkspb.LastAddedResponse{
		ID:    strconv.FormatInt(s.id, 10),
		Title: s.lastTitle,
		URL:   s.lastURL,
	}, nil
}

func startFlat(startCH chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", flatAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ser := grpc.NewServer(grpc.CustomCodec(flatbuffers.FlatbuffersCodec{}))
	bookmarks.RegisterBookmarksServiceServer(ser, &FlatServer{})
	fmt.Printf("Bookmarks FlatServer listening on %s\n", flatAddr)
	startCH <- struct{}{}
	if err := ser.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startGRPC(startCH chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ser := grpc.NewServer()
	bookmarkspb.RegisterBookmarksServiceServer(ser, &GRPCServer{})
	fmt.Printf("Bookmarks GRPCServer listening on %s\n", grpcAddr)
	startCH <- struct{}{}
	if err := ser.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	startCH := make(chan struct{})
	go startFlat(startCH, &wg)
	go startGRPC(startCH, &wg)
	<-startCH
	<-startCH
	wg.Wait()
}
