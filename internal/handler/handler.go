package handler

import (
	"encoding/json"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/response"
	"github.com/PereRohit/util/validator"
	"github.com/vatsal278/UserManagementService/internal/codes"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/logic"
	"github.com/vatsal278/UserManagementService/internal/model"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"io/ioutil"
	"net/http"
	"time"
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

func NewUserMgmtSvc(ds datasource.DataSourceI, jwtService jwtSvc.JWTService, msgQueue config.MsgQueue) UserMgmtSvcHandler {
	svc := &userMgmtSvc{
		logic: logic.NewUserMgmtSvcLogic(ds, jwtService, msgQueue),
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
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrReadingReqBody), nil)
		return
	}
	err = json.Unmarshal(body, &credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrUnmarshall), nil)
		return
	}
	err = validator.Validate(credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrValidate), nil)
		return
	}
	credential.RegistrationTimestamp, err = time.Parse("02-01-2006 15:04:05", credential.RegistrationDate)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrParseRegDate), nil)
		return
	}
	resp := svc.logic.Signup(credential)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
func (svc userMgmtSvc) Login(w http.ResponseWriter, r *http.Request) {
	var credential model.LoginCredentials
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrReadingReqBody), nil)
		return
	}
	err = json.Unmarshal(bytes, &credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrUnmarshall), nil)
		return
	}
	err = validator.Validate(credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrValidate), nil)
		return
	}
	resp := svc.logic.Login(w, credential)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)

}
