# University Management in Golang

## Starting docker containers
Run docker-compose up -d to start posrgres and pgadmin.
Default username is postgres and password `admin

## Connecting to pgadmin4
Use http://localhost:8080 and use username admin@admin.com and password admin

## To generate pb.go files
* go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
* go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
* export PATH="$PATH:$(go env GOPATH)/bin"
* Goto bin path and run "chmod +x protoc-gen-go"
* brew install protobuf
* protoc --version
* Go to the project directory and run "protoc --go-grpc_out=protoclient --go_out=protoclient university-management.proto"

## Pre-requisites
* brew install go@1.16
* go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate@latest
* go get github.com/shuLhan/go-bindata/...

## Migrations
* to create a new migration
`migrate create -ext sql -dir db/migrations -seq <name of the file>`

## To run the code
* Either go to the IDE and run the `server/main/main.go` and then `console_client/main/main.go`
  (OR)
* Open a terminal and run `go run server/main/main.go`
* Open another terminal and run `go run console_client/main/main.go`