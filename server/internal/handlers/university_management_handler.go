package handlers

import (
	"context"
	"encoding/json"
	"log"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"
)

type universityManagementServer struct {
	um.UniversityManagementServiceServer
	connectionManager connection.DatabaseConnectionManager
}

func (u *universityManagementServer) GetDepartment(ctx context.Context, request *um.GetDepartmentRequest) (*um.GetDepartmentResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var department *um.Department
	connection.GetSession().Select("id", "name").From("department").Where("id = ?", request.GetId()).LoadOne(&department)

	_, err = json.Marshal(&department)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}

	return &um.GetDepartmentResponse{Department: department}, nil
}

func (u *universityManagementServer) GetStudents(ctx context.Context, request *um.GetStudentRequest) (*um.GetStudentsResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var students []*um.Student

	var dep_id int32
	connection.GetSession().Select("id").From("departments").Where("name = ?", request.GetDepartmentName()).LoadOne(&dep_id)
	connection.GetSession().Select("roll", "name").From("student").Where("dep_id = ?", dep_id).Load(&students)

	return &um.GetStudentsResponse{Students: students}, nil
}

func (u *universityManagementServer) GetStudentDirectory(ctx context.Context, request *um.GetAllStudentRequest) (*um.GetAllStudentsResponse, error) {
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var students []*um.Student

	connection.GetSession().Select("roll", "name").From("student").Load(&students)

	return &um.GetAllStudentsResponse{Students: students}, nil
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
