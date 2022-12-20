package handler

import (
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/request"
	"github.com/PereRohit/util/response"
	"github.com/vatsal278/UserManagementService/internal/codes"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/logic"
	"github.com/vatsal278/UserManagementService/internal/model"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"github.com/vatsal278/UserManagementService/pkg/session"
	"net/http"
	"regexp"
	"time"
)

const UserManagementServiceName = "userManagementService"

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_handler.go --package=mock github.com/vatsal278/UserManagementService/internal/handler UserMgmtSvcHandler

type UserMgmtSvcHandler interface {
	HealthChecker
	SignUp(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	UserDetails(w http.ResponseWriter, r *http.Request)
	Activation(w http.ResponseWriter, r *http.Request)
}

type userMgmtSvc struct {
	logic logic.UserMgmtSvcLogicIer
}

func validatePass(pass string) *respModel.Response {
	if len(pass) < 8 {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Password must be 8 characters long",
			Data:    nil,
		}

	}
	done, err := regexp.MatchString("([a-z])+", pass)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to match password",
			Data:    nil,
		}
	}
	if !done {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Password must contain 1 lower case character",
			Data:    nil,
		}
	}
	done, err = regexp.MatchString("([A-Z])+", pass)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to match password",
			Data:    nil,
		}
	}
	if !done {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Password must contain 1 upper case character",
			Data:    nil,
		}
	}
	done, err = regexp.MatchString("([0-9])+", pass)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to match password",
			Data:    nil,
		}
	}
	if !done {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Password must contain 1 numeric character",
			Data:    nil,
		}
	}
	done, err = regexp.MatchString("([@_])+", pass)
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to match password",
			Data:    nil,
		}
	}
	if !done {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Password must contain 1 numeric character",
			Data:    nil,
		}
	}
	return &respModel.Response{Status: http.StatusOK}
}
func NewUserMgmtSvc(ds datasource.DataSourceI, jwtService jwtSvc.JWTService, msgQueue config.MsgQueue, cookie config.CookieStruct) UserMgmtSvcHandler {
	svc := &userMgmtSvc{
		logic: logic.NewUserMgmtSvcLogic(ds, jwtService, msgQueue, cookie),
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
	status, err := request.FromJson(r, &credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, status, err.Error(), nil)
		return
	}
	resp := validatePass(credential.Password)
	if resp.Status != http.StatusOK {
		response.ToJson(w, resp.Status, resp.Message, resp.Data)
		return
	}
	credential.RegistrationTimestamp, err = time.Parse("02-01-2006 15:04:05", credential.RegistrationDate)
	if err != nil {
		log.Error(err)
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrParseRegDate), nil)
		return
	}
	resp = svc.logic.Signup(credential)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
func (svc userMgmtSvc) Login(w http.ResponseWriter, r *http.Request) {
	var credential model.LoginCredentials
	status, err := request.FromJson(r, &credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, status, err.Error(), nil)
		return
	}
	resp := validatePass(credential.Password)
	if resp.Status != http.StatusOK {
		response.ToJson(w, resp.Status, resp.Message, resp.Data)
		return
	}
	resp = svc.logic.Login(w, credential)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
func (svc userMgmtSvc) Activation(w http.ResponseWriter, r *http.Request) {
	var data model.Activate
	status, err := request.FromJson(r, &data)
	if err != nil {
		log.Error(err)
		response.ToJson(w, status, err.Error(), nil)
		return
	}
	resp := svc.logic.Activate(data.UserId)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}

func (svc userMgmtSvc) UserDetails(w http.ResponseWriter, r *http.Request) {
	id := session.GetSession(r.Context())
	resp := svc.logic.UserData(id)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
