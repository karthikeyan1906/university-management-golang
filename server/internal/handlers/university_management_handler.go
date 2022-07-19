package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"university-management-golang/db/connection"
	um "university-management-golang/protoclient/university_management"

	"google.golang.org/protobuf/types/known/emptypb"
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
	log.Println("CaptureUserSignIn invoked")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var signInTime = req.GetSignInTime().AsTime()
	var formattedDate = signInTime.Format("2006-01-02")
	var studentId = req.GetRollnumber()
	var id int32

	var userActivity um.UserActivity
	uerr := connection.GetSession().QueryRow("SELECT id, studentid, signin, signout FROM user_activity WHERE studentid = $1 AND signin::date = $2", studentId, formattedDate).
		Scan(&userActivity.Id, &userActivity.Studentid, &userActivity.Signin, &userActivity.Signout)

	if uerr != nil && uerr != sql.ErrNoRows {
		log.Printf("Error while Capturing User Sign in - %v", uerr)
		return nil, uerr
	}

	// Student already logged in for the day
	if userActivity.GetId() != 0 {
		log.Printf("Old user with Id - %v\n", userActivity.GetId())
		return &um.SignInResponse{SignedInId: userActivity.GetId()}, nil
	}

	errs := connection.GetSession().QueryRow("INSERT INTO user_activity (studentid, signin) VALUES ($1, $2) RETURNING id", studentId, signInTime).Scan(&id)

	if errs != nil {
		log.Printf("Error while Capturing User Sign in - %v", errs)
		return nil, errs
	}

	return &um.SignInResponse{SignedInId: id}, nil
}

func (u *universityManagementServer) CaptureUserSignOut(ctx context.Context, req *um.SignOutRequest) (*emptypb.Empty, error) {
	log.Println("CaptureUserSignOut invoked")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}

	var studentId = req.GetRollnumber()
	var signedInId = req.GetSignedInId()
	var signOutTime = req.GetSignOutTime().AsTime()

	errs := connection.GetSession().QueryRow("UPDATE user_activity SET signout = $1 WHERE id = $2 AND studentid = $3", signOutTime, signedInId, studentId)

	return &emptypb.Empty{}, errs.Err()
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}
