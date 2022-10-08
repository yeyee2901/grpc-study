BASE_OUTPUT_DIR = ./gen


all: proto-go

proto-go:
	protoc  --go_out=${BASE_OUTPUT_DIR} --go_opt=paths=source_relative \
			--go-grpc_out=${BASE_OUTPUT_DIR} --go-grpc_opt=paths=source_relative \
			--go-grpc_opt=require_unimplemented_servers=false \
			proto/**/v1/*.proto
	protoc-go-inject-tag -input="./gen/proto/**/v1/*.pb.go"

clean:
	rm -rf ./gen/*
