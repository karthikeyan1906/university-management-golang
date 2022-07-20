package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"
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
	handleError(fmt.Sprintf("Error while creating DB connection: %+v", err), err)

	var id = request.GetId()
	log.Println(id)
	var department um.Department
	connection.GetSession().Select("id", "name").From("departments").Where("id = ?", request.GetId()).LoadOne(&department)

	_, err = json.Marshal(&department)
	handleError(fmt.Sprintf("Error while marshaling %+v", err), err)

	return &um.GetDepartmentResponse{Department: &um.Department{
		Id:   department.Id,
		Name: department.Name,
	}}, nil
}

func (u *universityManagementServer) GetStudents(ctx context.Context, req *um.GetStudentRequest) (*um.GetStudentResponse, error) {
	log.Println("Invoked Getstudent")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()
	handleError(fmt.Sprintf("Error while creating DB connection: %+v", err), err)

	var departmentName = req.GetDepartmentName()
	log.Printf("Input Dept name  is %v\n", departmentName)

	var students []um.Student

	_, sErr := connection.GetSession().Select("rollnumber", "students.name", "departmentid").From("students").
		Join("departments", "students.departmentid = departments.id").
		Where("? = '' OR departments.name = ?", departmentName, departmentName).
		Load(&students)
	handleError(fmt.Sprintf("Error while fetching student : %+v", sErr), sErr)

	studentsResp := formStudentResponse(students)

	return studentsResp, nil
}

func (u *universityManagementServer) CaptureUserSignIn(ctx context.Context, req *um.SignInRequest) (*um.SignInResponse, error) {
	log.Println("CaptureUserSignIn invoked")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()
	handleError(fmt.Sprintf("Error while creating DB connection: %+v", err), err)

	var signInTime = req.GetSignInTime().AsTime()
	var formattedDate = signInTime.Format("2006-01-02")
	var studentId = req.GetStudentId()
	var id int32

	if req.GetRollnumber() == 0 {
		go notifyLoginWithoutRollNumber(req)
	}

	userActivity, uerr := getUserActivityForSignIn(connection, studentId, formattedDate)
	if uerr != nil && uerr != sql.ErrNoRows {
		log.Printf("Error while Capturing User Sign in - %v", uerr)
		return nil, uerr
	}

	// Student already logged in for the day
	if userActivity.GetId() != 0 {
		log.Printf("Old user with Id - %v\n", userActivity.GetId())
		return &um.SignInResponse{SignedInId: userActivity.GetId()}, nil
	}

	go notifyNewLogin(req)

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
	handleError(fmt.Sprintf("Error while creating DB connection: %+v", err), err)

	var studentId = req.GetRollnumber()
	var signedInId = req.GetSignedInId()
	var signOutTime = req.GetSignOutTime().AsTime()

	userActivity, uerr := getUserActivityForSignOut(connection, signedInId, studentId)
	if uerr != nil && uerr != sql.ErrNoRows {
		log.Printf("Error while Capturing User Sign out - %v", uerr)
		return nil, uerr
	}

	checkAndNotifyEarlySignOut(&userActivity, signOutTime, req)

	errs := connection.GetSession().QueryRow("UPDATE user_activity SET signout = $1 WHERE id = $2 AND studentid = $3", signOutTime, signedInId, studentId)

	return &emptypb.Empty{}, errs.Err()
}

func (u *universityManagementServer) GetStaffs(ctx context.Context, req *um.GetStaffsRequest) (*um.GetStaffsResponse, error) {
	log.Println("Get Staffs invoked")

	connection, err := u.connectionManager.GetConnection()
	defer u.connectionManager.CloseConnection()
	handleError(fmt.Sprintf("Error while creating DB connection: %+v", err), err)

	var rollnumber = req.GetRollNumber()
	var staffs []um.Staff

	_, sErr := connection.GetSession().Select("staffs.id", "staffs.name").From("dept_staffs_mapping").
		Join("students", "dept_staffs_mapping.departmentid = students.departmentid").
		Join("staffs", "dept_staffs_mapping.staffid = staffs.id").
		Where("students.rollnumber = ?", rollnumber).
		Load(&staffs)
	handleError(fmt.Sprintf("Error while fetching staffs : %+v", sErr), sErr)

	staffsResp := formStaffResponse(staffs)

	return staffsResp, nil
}

func NewUniversityManagementHandler(connectionmanager connection.DatabaseConnectionManager) um.UniversityManagementServiceServer {
	return &universityManagementServer{
		connectionManager: connectionmanager,
	}
}

func formStudentResponse(students []um.Student) *um.GetStudentResponse {
	var studentsResp *um.GetStudentResponse = &um.GetStudentResponse{}
	for _, s := range students {
		student := um.Student{
			Rollnumber:   s.Rollnumber,
			Name:         s.Name,
			Departmentid: s.Departmentid,
		}

		studentsResp.Students = append(studentsResp.Students, &student)
	}

	return studentsResp
}

func getUserActivityForSignIn(connection connection.DatabaseConnect, studentId int32, formattedDate string) (um.UserActivity, error) {
	var id int32
	var studentIdMap int32
	var signIn sql.NullString
	var signOut sql.NullString

	uerr := connection.GetSession().QueryRow("SELECT id, studentid, signin, signout FROM user_activity WHERE studentid = $1 AND signin::date = $2", studentId, formattedDate).
		Scan(&id, &studentIdMap, &signIn, &signOut)

	return um.UserActivity{
		Id:        id,
		Studentid: studentIdMap,
		Signin:    signIn.String,
		Signout:   signOut.String,
	}, uerr
}

func getUserActivityForSignOut(connection connection.DatabaseConnect, signedInId int32, studentId int32) (um.UserActivity, error) {
	var id int32
	var studentIdMap int32
	var signIn sql.NullString
	var signOut sql.NullString

	uerr := connection.GetSession().QueryRow("SELECT id, studentid, signin, signout FROM user_activity WHERE id = $1 AND studentid = $2", signedInId, studentId).
		Scan(&id, &studentIdMap, &signIn, &signOut)

	return um.UserActivity{
		Id:        id,
		Studentid: studentIdMap,
		Signin:    signIn.String,
		Signout:   signOut.String,
	}, uerr
}

func formStaffResponse(staffs []um.Staff) *um.GetStaffsResponse {
	var staffsResp *um.GetStaffsResponse = &um.GetStaffsResponse{}
	for _, s := range staffs {
		staff := um.Staff{
			Id:   s.GetId(),
			Name: s.GetName(),
		}

		staffsResp.Staffs = append(staffsResp.Staffs, &staff)
	}

	return staffsResp
}

func checkAndNotifyEarlySignOut(userActivity *um.UserActivity, signOutTime time.Time, req *um.SignOutRequest) {
	layout := "2006-01-02T15:04:05.000000Z"
	t, cErr := time.Parse(layout, userActivity.Signin)
	handleError(fmt.Sprintf("Error while parsing signin time %+v", cErr), cErr)

	delta := signOutTime.Sub(t)
	log.Println(math.Round(delta.Hours()))

	if delta.Hours() < 8 {
		go notifyEarlySignOut(req, delta.Hours(), t)
	}
}

func notifyNewLogin(req *um.SignInRequest) {
	log.Printf("Student %s logged in at %v\n", req.GetStudentName(), req.GetSignInTime().AsTime())
}

func notifyLoginWithoutRollNumber(req *um.SignInRequest) {
	log.Printf("Student %s has logged in without rollnumber at %v\n", req.GetStudentName(), req.GetSignInTime().AsTime())
}

func notifyEarlySignOut(req *um.SignOutRequest, hours float64, signInTime time.Time) {
	log.Printf("Student %s has logged out in %.2f/8hrs from %v to %v\n", req.GetStudentName(), hours, signInTime, req.GetSignOutTime().AsTime())
}

func handleError(msg string, err error) {
	if err != nil {
		log.Fatalf(msg)
	}
}
