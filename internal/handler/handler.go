package handler

import (
	"fmt"
	"net/http"

	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"

	"github.com/vatsal278/UserManagementService/internal/logic"
	"github.com/vatsal278/UserManagementService/internal/model"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
)

const UserManagementServiceName = "userManagementService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/UserManagementService/internal/handler UserManagementServiceHandler

type UserManagementServiceHandler interface {
	HealthChecker
	Ping(w http.ResponseWriter, r *http.Request)
}

type userManagementService struct {
	logic logic.UserManagementServiceLogicIer
}

func NewUserManagementService(ds datasource.DataSource) UserManagementServiceHandler {
	svc := &userManagementService{
		logic: logic.NewUserManagementServiceLogic(ds),
	}
	AddHealthChecker(svc)
	return svc
}

func (svc userManagementService) HealthCheck() (svcName string, msg string, stat bool) {
	set := false
	defer func() {
		svcName = UserManagementServiceName
		if !set {
			msg = ""
			stat = true
		}
	}()
	stat = svc.logic.HealthCheck()
	set = true
	return
}

func (svc userManagementService) Ping(w http.ResponseWriter, r *http.Request) {
	req := &model.PingRequest{}

	suggestedCode, err := request.FromJson(r, req)
	if err != nil {
		response.ToJson(w, suggestedCode, fmt.Sprintf("FAILED: %s", err.Error()), nil)
		return
	}
	// call logic
	resp := svc.logic.Ping(req)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
	return
}
