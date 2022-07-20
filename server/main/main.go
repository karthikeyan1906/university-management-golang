package main

import (
	"log"
	"net"
	migrations "university-management-golang/db"
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
	err := migrations.MigrationsUp(username, password, host, dbPort, dbName, schema)
	if err != nil {
		log.Fatalf("Failed to migrate, err: %+v\n", err)
	}

	connectionmanager := &connection.DatabaseConnectionManagerImpl{
		DatabaseConfig: &connection.DBConfig{
			DbServer: host, DbPort: dbPort, DbUsername: username, DbPassword: password, DbNameSuffix: dbName, DbSchema: schema,
		},
		DatabaseConnection: nil,
	}

	//insertDepartmentSeedData(connectionmanager)
	//insertStudentSeedData(connectionmanager)
	//insertStaffsSeedData(connectionmanager)
	//insertDeptStaffsMappingSeedData(connectionmanager)

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

func insertDepartmentSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up department table")
	_, err = connection.GetSession().DeleteFrom("departments").Exec()
	if err != nil {
		log.Fatalf("Could not delete from departments table. Err: %+v", err)
	}

	log.Println("Inserting into department table")
	_, err = connection.GetSession().InsertInto("departments").Columns("id", "name").
		Values("1", "Computer Science").
		Values("2", "Information Technology").
		Values("3", "Electronic and Communication").
		Exec()

	if err != nil {
		log.Fatalf("Could not insert into departments table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}

func insertStudentSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up students table")
	_, err = connection.GetSession().DeleteFrom("students").Exec()
	if err != nil {
		log.Fatalf("Could not delete from students table. Err: %+v", err)
	}

	log.Println("Inserting into students table")
	_, err = connection.GetSession().InsertInto("students").Columns("rollnumber", "name", "departmentid").
		Values("1", "Virat Kholi", "2").
		Values("2", "Rohit Sharma", "2").
		Values("3", "Ravindra Jadeja", "2").
		Values("4", "Suryakumar Yadav", "3").
		Values("5", "Ravi Ashwin", "3").
		Exec()

	if err != nil {
		log.Fatalf("Could not insert into students table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}

func insertStaffsSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up staffs table")
	_, err = connection.GetSession().DeleteFrom("staffs").Exec()
	if err != nil {
		log.Fatalf("Could not delete from staffs table. Err: %+v", err)
	}

	log.Println("Inserting into staffs table")
	_, err = connection.GetSession().InsertInto("staffs").Columns("id", "name").
		Values("1", "Rahul Dravid").
		Values("2", "Sachin Tendulkar").
		Values("3", "Saurov Ganguly").
		Values("4", "VVS Laxman").
		Exec()

	if err != nil {
		log.Fatalf("Could not insert into staffs table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}

func insertDeptStaffsMappingSeedData(connectionManager connection.DatabaseConnectionManager) {
	connection, err := connectionManager.GetConnection()
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	log.Println("Cleaning up DeptStaffsMapping table")
	_, err = connection.GetSession().DeleteFrom("dept_staffs_mapping").Exec()
	if err != nil {
		log.Fatalf("Could not delete from DeptStaffsMapping table. Err: %+v", err)
	}

	log.Println("Inserting into DeptStaffsMapping table")
	_, err = connection.GetSession().InsertInto("dept_staffs_mapping").Columns("id", "departmentid", "staffid").
		Values("1", "2", "1").
		Values("2", "2", "2").
		Values("3", "2", "4").
		Values("4", "3", "3").
		Values("5", "3", "4").
		Exec()

	if err != nil {
		log.Fatalf("Could not insert into DeptStaffsMapping table. Err: %+v", err)
	}

	defer connectionManager.CloseConnection()
}
