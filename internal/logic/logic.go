package logic

import (
	"fmt"
	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/ggwhite/go-masker"
	"github.com/google/uuid"
	"github.com/vatsal278/UserManagementService/internal/codes"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/model"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/crypto"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"net/http"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/UserManagementService/internal/logic UserMgmtSvcLogicIer

type UserMgmtSvcLogicIer interface {
	HealthCheck() bool
	Signup(model.SignUpCredentials) *respModel.Response
	Login(http.ResponseWriter, model.LoginCredentials) *respModel.Response
	UserData(any) *respModel.Response
	Activate(id any) *respModel.Response
}

type userMgmtSvcLogic struct {
	DsSvc      datasource.DataSourceI
	jwtService jwtSvc.JWTService
	msgQueue   config.MsgQueue
	cookie     config.CookieStruct
}

func NewUserMgmtSvcLogic(ds datasource.DataSourceI, jwtService jwtSvc.JWTService, msgQueue config.MsgQueue, cookie config.CookieStruct) UserMgmtSvcLogicIer {
	return &userMgmtSvcLogic{
		DsSvc:      ds,
		jwtService: jwtService,
		msgQueue:   msgQueue,
		cookie:     cookie,
	}
}

func (l userMgmtSvcLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.DsSvc.HealthCheck()
}

func (l userMgmtSvcLogic) Signup(credential model.SignUpCredentials) *respModel.Response {
	result, err := l.DsSvc.Get(map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrCreatingAccount),
			Data:    nil,
		}
	}
	if len(result) != 0 {
		log.Error(codes.GetErr(codes.ErrEmailExists))
		return &respModel.Response{
			Status:  http.StatusBadRequest,
			Message: codes.GetErr(codes.ErrEmailExists),
			Data:    nil,
		}
	}
	newUser := model.User{
		Id:           uuid.New().String(),
		Email:        credential.Email,
		Name:         credential.Name,
		Company:      "",
		RegisteredOn: credential.RegistrationTimestamp,
	}

	hashedPassword, err := crypto.GeneratePasswordHash([]byte(credential.Password), []byte(newUser.Id))
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrHashPassword),
			Data:    nil,
		}
	}
	newUser.Password = hashedPassword

	err = l.DsSvc.Insert(newUser)
	if err != nil {
		log.Error(err.Error())
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrCreatingAccount),
			Data:    nil,
		}
	}
	go func(userId string, pubId string, channel string) {
		userID := fmt.Sprintf(`{"user_id":"%s"}`, userId)
		err := l.msgQueue.MsgBroker.PushMsg(userID, pubId, channel)
		if err != nil {
			log.Error(err)
			return
		}
		return
	}(newUser.Id, l.msgQueue.PubId, l.msgQueue.Channel)
	return &respModel.Response{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data:    "Account activation in progress",
	}
}

func (l userMgmtSvcLogic) Login(w http.ResponseWriter, credential model.LoginCredentials) *respModel.Response {
	result, err := l.DsSvc.Get(map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFetchingUser),
			Data:    nil,
		}
	}
	if len(result) == 0 {
		return &respModel.Response{
			Status:  http.StatusUnauthorized,
			Message: codes.GetErr(codes.AccNotFound),
			Data:    nil,
		}
	}
	log.Info(result[0].Id)
	hashedPassword, err := crypto.GeneratePasswordHash([]byte(credential.Password), []byte(result[0].Id))
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrHashPassword),
			Data:    nil,
		}
	}
	result, err = l.DsSvc.Get(map[string]interface{}{"email": credential.Email, "password": hashedPassword})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrFetchingUser),
			Data:    nil,
		}
	}
	if len(result) == 0 {
		return &respModel.Response{
			Status:  http.StatusUnauthorized,
			Message: codes.GetErr(codes.InvaliCredentials),
			Data:    nil,
		}
	}
	if result[0].Active != true {
		return &respModel.Response{
			Status:  http.StatusAccepted,
			Message: codes.GetErr(codes.ErrLogging),
			Data:    codes.GetErr(codes.AccActivationInProcess),
		}
	}
	id := result[0].Id
	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrDuration),
			Data:    nil,
		}
	}
	jwtToken, err := l.jwtService.GenerateToken(jwt.SigningMethodHS256, id, l.cookie.Expiry)
	if err != nil {
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrGenerateJwt),
			Data:    nil,
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		MaxAge:   int(l.cookie.Expiry),
		Path:     l.cookie.Path,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(l.cookie.Expiry),
	})

	newActiveDvc := result[0].ActiveDevices + 1
	err = l.DsSvc.Update(map[string]interface{}{"active_devices": newActiveDvc}, map[string]interface{}{"email": credential.Email})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: codes.GetErr(codes.ErrLogging),
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: codes.GetErr(codes.Success),
		Data:    nil,
	}
}

func (l userMgmtSvcLogic) Activate(id any) *respModel.Response {

	err := l.DsSvc.Update(map[string]interface{}{"active": true}, map[string]interface{}{"user_id": id})
	if err != nil {
		log.Error(err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "cant activate the account",
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    nil,
	}
}

func (l userMgmtSvcLogic) UserData(id any) *respModel.Response {
	var userDetails model.UserDetails
	i, ok := id.(string)
	if !ok {
		log.Error("cant assert id")
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "cant assert id",
			Data:    nil,
		}
	}
	user, err := l.DsSvc.Get(map[string]interface{}{"user_id": i})
	if err != nil {
		log.Error("cant fetch user from db")
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "cant fetch user from db",
			Data:    nil,
		}
	}
	userDetails.Name = user[len(user)-1].Name

	userDetails.Email = masker.Email(user[len(user)-1].Email)
	userDetails.Company = user[len(user)-1].Company
	userDetails.LastLogin = user[len(user)-1].UpdatedOn

	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		Data:    userDetails,
	}
}
