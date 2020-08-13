

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative api.proto api_options.proto
	awk -f extract.awk api.proto > messages.go
