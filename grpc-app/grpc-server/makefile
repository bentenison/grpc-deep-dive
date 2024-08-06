grpc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto

client:
	go run ../../grpc-client/client.go

server:
	go run main.go