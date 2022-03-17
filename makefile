grpc_db_set:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/grpc_db/grpc_db.proto

ftw:
	go get google.golang.org/protobuf
	go get google.golang.org/grpc
	go get github.com/aws/aws-sdk-go/aws

run_client:
	go run cmd/client.go

run_server_debug:
	go run -race cmd/server.go -v

run_server:
	go run cmd/server.go

run_multiple_servers:
	./cmd/run_servers.sh

run_multiple_servers_debug:
	./cmd/run_servers.sh -v

shutdown_multiple_servers:
	./cmd/shutdown_servers.sh

run_stream:
	go run cmd/stream.go

run_sensors_simulator:
	go run cmd/sensor_simulator.go
