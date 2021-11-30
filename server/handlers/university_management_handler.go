package handlers

import (
	"context"
	"log"
	um "university-management-golang/protoclient/university_management"
)

type universityManagementServer struct {
	um.UniversityManagementServiceServer
}

func (u *universityManagementServer) GetDepartment(ctx context.Context, request *um.GetDepartmentRequest) (*um.GetDepartmentResponse, error) {
	log.Println("Hello World!")
	return nil, nil
}

func NewUniversityManagementHandler() um.UniversityManagementServiceServer {
	return &universityManagementServer{}
}