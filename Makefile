.PHONY: generate_fbs compile compile_bookmarks_client compile_bookmarks_server all

all: generate_fbs generate_proto compile_bookmarks_client compile_bookmarks_server

generate_fbs:
	flatc --go --grpc bookmarks.fbs

generate_proto:
	protoc bookmarks.proto --go_out 'bookmarkspb' --go_opt 'paths=source_relative' --go-grpc_out 'bookmarkspb' --go-grpc_opt 'paths=source_relative' #--go_out=plugins=grpc:bookmarkspb

compile_bookmarks_client:
	go build -o client ./bookmarks-client

compile_bookmarks_server:
	go build -o server ./bookmarks-server
