.PHONY: generate_fbs compile compile_bookmarks_client compile_bookmarks_server all

all: generate_fbs compile_bookmarks_client compile_bookmarks_server

generate_fbs:
	flatc --go --grpc bookmarks.fbs

compile_bookmarks_client:
	go build -o client ./bookmarks-client

compile_bookmarks_server:
	go build -o server ./bookmarks-server
