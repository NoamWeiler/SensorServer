dbset:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/mutual_db/mutual_db.proto

ftw:
	go get google.golang.org/protobuf
	go get google.golang.org/grpc
	go get google.golang.org/grpc
	go get github.com/aws/aws-sdk-go/aws

run_client:
	go run pkg/client/main.go pkg/client/client_menu.go

run_server:
	go run pkg/server/*.go -v

run_stream:
	go run pkg/sensor/main.go
