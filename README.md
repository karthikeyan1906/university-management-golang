# University Management in Golang

## Starting docker containers
Run `docker-compose up -d` to start posrgres and pgadmin.
Default username is `postgres` and password `admin

## Connecting to pgadmin4
Use http://localhost:8080 and use username `admin@admin.com` and password `admin`

## To generate pb.go files
* go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
* go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
* export PATH="$PATH:$(go env GOPATH)/bin"
* Goto bin path and do a chmod +x on protoc-gen-go
* brew install protobuf
* protoc --version 
* protoc --go-grpc_out=protoclient --go_out=protoclient university-management.proto
