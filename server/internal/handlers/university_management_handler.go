package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
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

// func (u *universityManagementServer) GetStudentDirectory(ctx context.Context, request *um.GetAllStudentRequest) (*um.GetAllStudentsResponse, error){
// 	connection, err := u.connectionManager.GetConnection()
// 	defer u.connectionManager.CloseConnection()

// 	if err != nil {
// 		log.Fatalf("Error: %+v", err)
// 	}

// 	var students []*um.Student

// 	connection.GetSession().Select("roll", "name").From("student").Load(&students)

// 	return &um.GetAllStudentsResponse{Students: students},nil
// }

func (u *universityManagementServer) Notify(ctx context.Context, request *um.GetNotifyRequest) (*um.GetNotifyResponse, error) {
	c := make(chan string)

	id := request.GetId()
	log.Println("Strting Go routine")
	go u.waitForStudentToLogin(id, c)

	message := <-c

	return &um.GetNotifyResponse{Message: message}, nil
}

func (u *universityManagementServer) waitForStudentToLogin(id int32, c chan string) {
	log.Println("In Go routine")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	for {

		time.Sleep(1 * time.Second)
		loginTime := u.getLoginTime(id,connection)

		if loginTime != nil {
			c <- "Hello I am loggedIn"
			return
		}
	}
}

func (u *universityManagementServer) getLoginTime(id int32, connection connection.DatabaseConnect) *time.Time {

	var timelogin *time.Time

	connection.GetSession().Select("loginTime").From("attendance").Where("id = ?", id).LoadOne(&timelogin)

	log.Printf("In get time %v", timelogin)
	return timelogin
}
func difference(a, b []int32) (diff []int32) {
	m := make(map[int32]bool)

	for _, item := range b {
			m[item] = true
	}

	for _, item := range a {
			if _, ok := m[item]; !ok {
					diff = append(diff, item)
			}
	}
	return
}
func ( u *universityManagementServer) GetAttendance(request *um.GetAttendanceRequest, stream um.UniversityManagementService_GetAttendanceServer) error{
	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	log.Println("In server streaming")
	var loggedIds []int32
	var sentIds []int32
	for{
		time.Sleep(5*time.Second)

		connection.GetSession().Select("id").From("attendance").Load(&loggedIds)
		//log.Println("from database",loggedIds)
		if len(loggedIds)>0{
			log.Println("Sending stream")
			toSend:=difference(loggedIds,sentIds)
			stream.Send(&um.GetAttendanceResponse{Ids:toSend})
			sentIds=append(sentIds,toSend...)

		}
	}
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
