package main

import (
	"log"
	"net"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"
	"university-management-golang/server/internal/handlers"

	"google.golang.org/grpc"
)

const port = "2345"

//db
const (
	username = "postgres"
	password = "admin"
	host     = "localhost"
	dbPort   = "5436"
	dbName   = "postgres"
	schema   = "public"
)

func main() {
	// err := migrations.MigrationsUp(username, password, host, dbPort, dbName, schema)
	// if err != nil {
	// 	log.Fatalf("Failed to migrate, err: %+v\n", err)
	// }

	connectionmanager := &connection.DatabaseConnectionManagerImpl{
		&connection.DBConfig{
			host, dbPort, username, password, dbName, schema,
		},
		nil,
	}

	//insertSeedData(connectionmanager)

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen to port: %s, err: %+v\n", port, err)
	}
	log.Printf("Starting to listen on port: %s\n", port)

	um.RegisterUniversityManagementServiceServer(grpcServer, handlers.NewUniversityManagementHandler(connectionmanager))
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to start GRPC Server: %+v\n", err)
	}
}

func insertSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up department table")
	_, err = connection.GetSession().DeleteFrom("department").Exec()
	if err != nil {
		log.Fatalf("Could not delete from department table. Err: %+v", err)
	}

	log.Println("Inserting into department table")
	_, err = connection.GetSession().InsertInto("department").Columns("name").
		Values("Computer Science").Exec()

	if err != nil {
		log.Fatalf("Could not insert into department table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
