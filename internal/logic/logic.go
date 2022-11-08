package logic

import (
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/UserManagementService/internal/logic UserManagementServiceLogicIer

type UserMgmtSvcLogicIer interface {
	HealthCheck() bool
	//Signup(model.SignUpCredentials) *respModel.Response
	//Login(http.ResponseWriter, model.LoginCredentials) *respModel.Response
	//UserData(any) *respModel.Response
	//Activate(id any) *respModel.Response
}

type userMgmtSvcLogic struct {
	DsSvc datasource.DataSourceI
	//loginService jwtSvc.LoginService
	//jwtService   jwtSvc.JWTService
}

func NewUserMgmtSvcLogic(ds datasource.DataSourceI) UserMgmtSvcLogicIer {
	return &userMgmtSvcLogic{
		DsSvc: ds,
		//loginService: loginSvc,
		//jwtService:   jwtService,
	}
}

func (l userMgmtSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.DsSvc.HealthCheck()
}
