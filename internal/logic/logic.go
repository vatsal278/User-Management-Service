package logic

import (
	"fmt"
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/vatsal278/UserManagementService/internal/model"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/bCrypt"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"github.com/vatsal278/UserManagementService/internal/repo/helpers"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"net/http"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/UserManagementService/internal/logic UserManagementServiceLogicIer

type UserMgmtSvcLogicIer interface {
	HealthCheck() bool
	Signup(model.SignUpCredentials) *respModel.Response
	Login(http.ResponseWriter, model.LoginCredentials) *respModel.Response
	//UserData(any) *respModel.Response
	//Activate(id any) *respModel.Response
}

type userMgmtSvcLogic struct {
	DsSvc        datasource.DataSourceI
	loginService helpers.LoginService
	jwtService   jwtSvc.JWTService
	msgQueue     sdk.MsgBrokerSvcI
}

func NewUserMgmtSvcLogic(ds datasource.DataSourceI, loginSvc helpers.LoginService, jwtService jwtSvc.JWTService, msgQueue sdk.MsgBrokerSvcI) UserMgmtSvcLogicIer {
	return &userMgmtSvcLogic{
		DsSvc:        ds,
		loginService: loginSvc,
		jwtService:   jwtService,
		msgQueue:     msgQueue,
	}
}

func (l userMgmtSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.DsSvc.HealthCheck()
}

func (l userMgmtSvcLogic) Signup(credential model.SignUpCredentials) *respModel.Response {
	if credential.Email == "" || credential.Password == "" || credential.Name == "" {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Provide required details",
			Data:    nil,
		}
	}
	result, err := l.DsSvc.Get(map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Problem creating account",
			Data:    nil,
		}
	}
	if len(result) != 0 {
		log.Error("Email is already in use")
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Email is already in use",
			Data:    nil,
		}
	}

	newUser := model.User{
		Id:           uuid.New().String(),
		Email:        credential.Email,
		Password:     credential.Password,
		Name:         credential.Name,
		Company:      "perennial",
		RegisteredOn: time.Now(),
	}
	salt := bCrypt.GeneratePasswordHash([]byte(newUser.Password))
	newUser.Salt = salt
	err = l.DsSvc.Insert(newUser)
	log.Info(newUser.Id)
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Problem creating an account",
			Data:    nil,
		}
	}
	id, err := l.msgQueue.RegisterPub("New Account Activation Request Channel")
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Problem registering to account management service",
			Data:    nil,
		}
	}
	userID := fmt.Sprintf(`{"user_id":"%s"}`, newUser.Id)
	err = l.msgQueue.PushMsg(userID, id, "New Account Activation Request")
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Problem notifying account management service",
			Data:    nil,
		}
	}

	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data:    nil,
	}
}

func (l userMgmtSvcLogic) Login(w http.ResponseWriter, credential model.LoginCredentials) *respModel.Response {
	if credential.Email == "" || credential.Password == "" {
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: "Provide required details",
			Data:    nil,
		}
	}
	result, err := l.DsSvc.Get(map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Problem logging into your account",
			Data:    nil,
		}
	}
	if result[0].Email == "" {
		return &respModel.Response{
			Status:  http.StatusUnauthorized,
			Message: "User account was not found",
			Data:    nil,
		}
	}
	if result[0].Active != true {
		return &respModel.Response{
			Status:  http.StatusAccepted,
			Message: "Problem logging into your account",
			Data:    "Account activation in progress",
		}
	}

	hashedPassword := []byte(result[0].Salt)
	// Get the password provided in the request.body
	password := []byte(credential.Password)

	err = bCrypt.PasswordCompare(password, hashedPassword)

	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Invalid user credentials",
			Data:    nil,
		}
	}
	id := result[0].Id
	jwtToken, err := l.jwtService.GenerateToken(jwt.SigningMethodHS256, id, 6)
	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Unable to generate jwt token",
			Data:    nil,
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: jwtToken,
	})
	newActiveDvc := result[0].ActiveDevices + 1
	err = l.DsSvc.Update(map[string]interface{}{"active_devices": newActiveDvc}, map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "Problem logging into your account",
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    nil,
	}
}
