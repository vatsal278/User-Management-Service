package handler

import (
	"encoding/json"
	"github.com/PereRohit/util/response"
	"github.com/vatsal278/UserManagementService/internal/logic"
	"github.com/vatsal278/UserManagementService/internal/model"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"github.com/vatsal278/UserManagementService/internal/repo/helpers"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"io/ioutil"
	"net/http"
)

const UserManagementServiceName = "userManagementService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/UserManagementService/internal/handler UserManagementServiceHandler

type UserMgmtSvcHandler interface {
	HealthChecker
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	//UserDetails(w http.ResponseWriter, r *http.Request)
	//Activation(w http.ResponseWriter, r *http.Request)
}

type userMgmtSvc struct {
	logic logic.UserMgmtSvcLogicIer
}

func NewUserMgmtSvc(ds datasource.DataSourceI, loginService helpers.LoginService, jwtService jwtSvc.JWTService, msgQueue sdk.MsgBrokerSvcI) UserMgmtSvcHandler {
	svc := &userMgmtSvc{
		logic: logic.NewUserMgmtSvcLogic(ds, loginService, jwtService, msgQueue),
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
func (svc userMgmtSvc) SignUp(w http.ResponseWriter, r *http.Request) {
	var credential model.SignUpCredentials

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &credential)
	if err != nil {
		return
	}
	resp := svc.logic.Signup(credential)

	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
func (svc userMgmtSvc) Login(w http.ResponseWriter, r *http.Request) {
	var credential model.LoginCredentials
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &credential)
	if err != nil {
		return
	}
	resp := svc.logic.Login(w, credential)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)

}
