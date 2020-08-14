
.PHONY: all
all: proto message
	go build

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative api.proto api_options.proto

message:
	awk -f extract.awk api.proto | gofmt > message.go
