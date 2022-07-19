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

	var departmentID int32 = 1
	departmentResponse, err := univClient.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	if err != nil {
		log.Fatalf("Error occured while fetching department for id %d,err: %+v \n", departmentID, err)
	}
	log.Println(departmentResponse)

	var departmentName string = "Information Technology"
	resp, errr := univClient.GetStudents(context.TODO(), &university_management.GetStudentRequest{DepartmentName: departmentName})
	if errr != nil {
		log.Fatalf("Error occured while fetching students for id %s,err: %+v \n", departmentName, errr)
	}
	log.Println(resp)

	var studentId int32 = 2

	signInResp, e := univClient.CaptureUserSignIn(context.TODO(), &university_management.SignInRequest{
		Rollnumber: studentId,
		SignInTime: timestamppb.Now(),
	})

	if e != nil {
		log.Fatalf("Error occured while adding sign in time for student id %d, err : %v \n", studentId, e)
	} else {
		log.Printf("Captured User sign in time with Id - %d", signInResp.GetSignedInId())
	}

	_, errs := univClient.CaptureUserSignOut(context.TODO(), &university_management.SignOutRequest{
		Rollnumber:  studentId,
		SignOutTime: timestamppb.Now(),
		SignedInId:  signInResp.GetSignedInId(),
	})

	if errs != nil {
		log.Fatalf("Error occured while adding sign out time for student id %d, err : %v \n", studentId, errs)
	} else {
		log.Printf("Captured User sign out time for Id - %d", signInResp.GetSignedInId())
	}

}
