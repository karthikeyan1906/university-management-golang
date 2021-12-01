package main

import (
	"context"
	"fmt"
	"log"
	"university-management-golang/protoclient/university_management"

	"google.golang.org/grpc"
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
	client := university_management.NewUniversityManagementServiceClient(conn)
	// var departmentID int32 = 1
	// departmentResponse, err := client.GetDepartment(context.TODO(), &university_management.GetDepartmentRequest{Id: departmentID})
	// if err != nil {
	// 	log.Fatalf("Error occured while fetching department for id %d,err: %+v", departmentID, err)
	// }
	// log.Println(departmentResponse)

	// var departName string = "CS"
	// studentResponse, err := client.GetStudents(context.TODO(), &university_management.GetStudentRequest{DepartmentName: departName})

	// if err != nil {
	// 	log.Fatalf("Error occured while fetching students for department %s,err: %+v", departName, err)
	// }

	// log.Println(studentResponse)

	log.Println("******All Student Api******")

	studentAllResponse, err := client.GetStudentDirectory(context.TODO(), &university_management.GetAllStudentRequest{})

	if err != nil {
		log.Fatalf("Error occured while fetching students,err: %+v", err)
	}

	log.Println(studentAllResponse)
}
