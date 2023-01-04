package handler

import (
	"errors"
	"github.com/PereRohit/util/log"
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

func validatePass(pass string) error {
	done, err := regexp.MatchString("([a-z])+", pass)
	if err != nil {
		return errors.New(codes.GetErr(codes.ErrPassRegex))
	}
	if !done {
		return errors.New(codes.GetErr(codes.ErrPassLowerCase))
	}
	done, err = regexp.MatchString("([A-Z])+", pass)
	if err != nil {
		return errors.New(codes.GetErr(codes.ErrPassRegex))
	}
	if !done {
		return errors.New(codes.GetErr(codes.ErrPassUpperCase))
	}
	done, err = regexp.MatchString("([0-9])+", pass)
	if err != nil {
		return errors.New(codes.GetErr(codes.ErrPassRegex))
	}
	if !done {
		return errors.New(codes.GetErr(codes.ErrPassNumeric))
	}
	done, err = regexp.MatchString("([@.,$?])+", pass)
	if err != nil {
		return errors.New(codes.GetErr(codes.ErrPassRegex))
	}
	if !done {
		return errors.New(codes.GetErr(codes.ErrPassSpecial))
	}
	return nil
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
	err = validatePass(credential.Password)
	if err != nil {
		log.Error(err)
		if err.Error() != codes.GetErr(codes.ErrPassRegex) {
			response.ToJson(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.ToJson(w, http.StatusInternalServerError, err.Error(), nil)
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
	status, err := request.FromJson(r, &credential)
	if err != nil {
		log.Error(err)
		response.ToJson(w, status, err.Error(), nil)
		return
	}
	err = validatePass(credential.Password)
	if err != nil {
		log.Error(err)
		if err.Error() != codes.GetErr(codes.ErrPassRegex) {
			response.ToJson(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.ToJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resp := svc.logic.Login(w, credential)
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
	idStr, ok := id.(string)
	if !ok {
		response.ToJson(w, http.StatusBadRequest, codes.GetErr(codes.ErrAssertUserid), nil)
		return
	}
	resp := svc.logic.UserData(idStr)
	response.ToJson(w, resp.Status, resp.Message, resp.Data)
}
