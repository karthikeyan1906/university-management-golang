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

	var id = request.GetId()
	log.Println(id)
	var department um.Department
	connection.GetSession().Select("id", "name").From("department").Where("id = ?", request.GetId()).LoadOne(&department)

	_, err = json.Marshal(&department)
	if err != nil {
		log.Fatalf("Error while marshaling %+v", err)
	}

	return &um.GetDepartmentResponse{Department: &um.Department{
		Id:   department.Id,
		Name: department.Name,
	}}, nil
}

func (u *universityManagementServer) GetStudents(ctx context.Context, req *um.GetStudentRequest) (*um.GetStudentResponse, error) {
	log.Println("Invoked Getstudent")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var departmentName = req.GetDepartmentName()
	log.Printf("Input Dept name  is %v\n", departmentName)

	var students []um.Student

	_, errors := connection.GetSession().Select("rollnumber", "students.name", "departmentid").From("students").Join("departments", "students.departmentid = departments.id").Where("departments.name = ?", departmentName).Load(&students)

	if errors != nil {
		log.Fatalf("Error: %+v", errors)
	}

	var studentsResp *um.GetStudentResponse = &um.GetStudentResponse{}

	for _, s := range students {
		student := um.Student{
			Rollnumber:   s.Rollnumber,
			Name:         s.Name,
			Departmentid: s.Departmentid,
		}

		studentsResp.Students = append(studentsResp.Students, &student)
	}

	return studentsResp, nil
}

func (u *universityManagementServer) CaptureUserSignIn(ctx context.Context, req *um.SignInRequest) (*um.SignInResponse, error) {
	log.Println("AddUserSignIn invoked")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var signInTime = req.GetSignInTime().AsTime()
	var studentId = req.GetRollnumber()

	// var id int

	res, e := connection.GetSession().InsertInto("user_activity").Columns("studentid", "signin").Values(studentId, signInTime).Returning("id").Exec()

	// var id int
	// var query string = fmt.Sprintf("INSERT INTO user_activity (studentid, signin) VALUES ( %d, %s) RETURNING id;", studentId, signInTime)
	// errs := connection.GetSession().QueryRow(query).Scan(&id)
	// if errs != nil {
	// 	log.Fatalln(errs)
	// }

	log.Println(res.LastInsertId())

	return &um.SignInResponse{SignInId: 1}, e
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
