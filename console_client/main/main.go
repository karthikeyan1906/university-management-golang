package main

import (
	"context"
	"fmt"
	"log"
	"university-management-golang/protoclient/university_management"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	host = "localhost"
	port = "2345"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error occured %+v", err)
	}
	univClient := university_management.NewUniversityManagementServiceClient(conn)

	var departmentID int32 = 2
	departmentResponse, dErr := univClient.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	if dErr != nil {
		log.Fatalf("Error occured while fetching department for id %d, err: %+v \n", departmentID, dErr)
	}
	log.Println(departmentResponse)

	var departmentName string = "Information Technology"
	studResp, sErr := univClient.GetStudents(context.TODO(), &university_management.GetStudentRequest{DepartmentName: departmentName})
	if sErr != nil {
		log.Fatalf("Error occured while fetching students for id %s, err: %+v \n", departmentName, sErr)
	}
	log.Println(studResp)

	var studentId int32 = 3
	signInResp, siErr := univClient.CaptureUserSignIn(context.TODO(), &university_management.SignInRequest{
		Rollnumber:  studentId,
		SignInTime:  timestamppb.Now(),
		StudentName: "Test2",
		StudentId:   studentId,
	})

	if siErr != nil {
		log.Fatalf("Error occured while adding sign in time for student id %d, err : %v \n", studentId, siErr)
	} else {
		log.Printf("Captured User sign in time with Id - %d", signInResp.GetSignedInId())
	}

	_, soErr := univClient.CaptureUserSignOut(context.TODO(), &university_management.SignOutRequest{
		Rollnumber:  studentId,
		SignOutTime: timestamppb.Now(),
		SignedInId:  signInResp.GetSignedInId(),
	})

	if soErr != nil {
		log.Fatalf("Error occured while adding sign out time for student id %d, err : %v \n", studentId, soErr)
	} else {
		log.Printf("Captured User sign out time for Id - %d", signInResp.GetSignedInId())
	}

	logInResp, sgiErr := univClient.CaptureUserSignIn(context.TODO(), &university_management.SignInRequest{
		SignInTime:  timestamppb.Now(),
		StudentName: "Test2",
		StudentId:   studentId,
	})

	if sgiErr != nil {
		log.Fatalf("Error occured while adding sign in time for student id %d, err : %v \n", studentId, sgiErr)
	} else {
		log.Printf("Captured User sign in time with Id - %d", logInResp.GetSignedInId())
	}
}
