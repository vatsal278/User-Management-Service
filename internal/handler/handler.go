package handler

import (
	"github.com/vatsal278/UserManagementService/internal/logic"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
)

const UserManagementServiceName = "userManagementService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/UserManagementService/internal/handler UserManagementServiceHandler

type UserMgmtSvcHandler interface {
	HealthChecker
	//SignUp(w http.ResponseWriter, r *http.Request)
	//Login(w http.ResponseWriter, r *http.Request)
	//UserDetails(w http.ResponseWriter, r *http.Request)
	//Activation(w http.ResponseWriter, r *http.Request)
}

type userMgmtSvc struct {
	logic logic.UserMgmtSvcLogicIer
}

func NewUserMgmtSvc(ds datasource.DataSourceI) UserMgmtSvcHandler {
	svc := &userMgmtSvc{
		logic: logic.NewUserMgmtSvcLogic(ds),
	}
	AddHealthChecker(svc)
	return svc
}

func (svc userMgmtSvc) HealthCheck() (svcName string, msg string, stat bool) {
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
