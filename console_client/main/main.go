package main

import (
	"context"
	"fmt"
	"log"
	"time"
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

	// List department
	var departmentID int32 = 2
	departmentResponse, dErr := univClient.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	if dErr != nil {
		log.Fatalf("Error occured while fetching department for id %d, err: %+v \n", departmentID, dErr)
	}
	log.Println(departmentResponse)

	// List students from a department
	var departmentName string = "Information Technology"
	studResp, sErr := univClient.GetStudents(context.TODO(), &university_management.GetStudentRequest{DepartmentName: departmentName})
	if sErr != nil {
		log.Fatalf("Error occured while fetching students for id %s, err: %+v \n", departmentName, sErr)
	}
	log.Println(studResp)

	//Capture Student Sign in time
	var studentId int32 = 2
	signInResp, siErr := univClient.CaptureUserSignIn(context.TODO(), &university_management.SignInRequest{
		Rollnumber:  studentId,
		SignInTime:  timestamppb.Now(),
		StudentName: "Rohit Sharma",
		StudentId:   studentId,
	})

	if siErr != nil {
		log.Fatalf("Error occured while adding sign in time for student id %d, err : %v \n", studentId, siErr)
	} else {
		log.Printf("Captured User sign in time with Id - %d", signInResp.GetSignedInId())
	}

	//Capture Student Sign out time
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

	//Capture Student Sign in time along with notification for sign in without rollnumber
	logInResp, sgiErr := univClient.CaptureUserSignIn(context.TODO(), &university_management.SignInRequest{
		SignInTime:  timestamppb.Now(),
		StudentName: "Rohit Sharma",
		StudentId:   studentId,
	})

	if sgiErr != nil {
		log.Fatalf("Error occured while adding sign in time for student id %d, err : %v \n", studentId, sgiErr)
	} else {
		log.Printf("Captured User sign in time without rollnumber with Id - %d", logInResp.GetSignedInId())
	}

	//Capture Student Sign out time with notification for early sign out (within 8 hours from sign in time)
	_, sonErr := univClient.CaptureUserSignOut(context.TODO(), &university_management.SignOutRequest{
		Rollnumber:  studentId,
		SignOutTime: timestamppb.Now(),
		SignedInId:  logInResp.GetSignedInId(),
		StudentName: "Rohit Sharma",
	})

	if sonErr != nil {
		log.Fatalf("Error occured while adding sign out time for student id %d, err : %v \n", 2, sonErr)
	} else {
		log.Printf("Captured User early sign out time for Id - %d", logInResp.GetSignedInId())
	}

	//Capture Student Sign out time after 8 hours from sign in time without notification
	_, soNotiErr := univClient.CaptureUserSignOut(context.TODO(), &university_management.SignOutRequest{
		Rollnumber:  studentId,
		SignOutTime: timestamppb.New(time.Now().Add(time.Hour * time.Duration(9))),
		SignedInId:  logInResp.GetSignedInId(),
		StudentName: "Rohit Sharma",
	})

	if soNotiErr != nil {
		log.Fatalf("Error occured while adding sign out time for student id %d, err : %v \n", 2, soNotiErr)
	} else {
		log.Printf("Captured User sign out time after 8 hours for Id - %d", logInResp.GetSignedInId())
	}

	// List students directory
	var empty string = ""
	stuResp, stuErr := univClient.GetStudents(context.TODO(), &university_management.GetStudentRequest{DepartmentName: empty})
	if stuErr != nil {
		log.Fatalf("Error occured while fetching students directory for id %s, err: %+v \n", departmentName, stuErr)
	}
	log.Println(stuResp)
}
