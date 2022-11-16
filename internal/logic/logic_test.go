package logic

import (
	"errors"
	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/vatsal278/UserManagementService/internal/codes"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/model"
	"github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"github.com/vatsal278/UserManagementService/pkg/mock"
)

func Test_userManagementServiceLogic_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want  bool
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)

				mockDs.EXPECT().HealthCheck().Times(1).
					Return(true)

				return mockDs, nil, config.MsgQueue{}, config.CookieStruct{}
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewUserMgmtSvcLogic(tt.setup())

			got := rec.HealthCheck()

			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func Test_userManagementServiceLogic_SignUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials model.SignUpCredentials
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct)
		want        func(*respModel.Response)
	}{
		{
			name: "Success",
			credentials: model.SignUpCredentials{
				Name:                  "Vatsal",
				Email:                 "vatsal@gmail.com",
				Password:              "Abcde@12345",
				RegistrationTimestamp: time.Now(),
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return([]model.User{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    codes.GetErr(codes.AccActivationInProcess),
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Success :: Push msg failure",
			credentials: model.SignUpCredentials{
				Name:                  "Vatsal",
				Email:                 "vatsal@gmail.com",
				Password:              "Abcde@12345",
				RegistrationTimestamp: time.Now(),
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return([]model.User{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    codes.GetErr(codes.AccActivationInProcess),
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure::Get from db err",
			credentials: model.SignUpCredentials{
				Name:                  "Vatsal",
				Email:                 "vatsal@gmail.com",
				Password:              "Abcde@12345",
				RegistrationTimestamp: time.Now(),
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue, config.CookieStruct) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New(""))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}, config.CookieStruct{}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrCreatingAccount),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure:: Email already exists",
			credentials: model.SignUpCredentials{
				Name:                  "Vatsal",
				Email:                 "vatsal@gmail.com",
				Password:              "Abcde@12345",
				RegistrationTimestamp: time.Now(),
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com"})
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrEmailExists),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure:: Error Inserting in db",
			credentials: model.SignUpCredentials{
				Name:                  "Vatsal",
				Email:                 "vatsal@gmail.com",
				Password:              "Abcde@12345",
				RegistrationTimestamp: time.Now(),
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return([]model.User{}, nil)
				mockDs.EXPECT().Insert(gomock.Any()).Times(1).Return(errors.New(""))
				return mockDs, nil, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9091")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrCreatingAccount),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewUserMgmtSvcLogic(tt.setup())

			got := rec.Signup(tt.credentials)

			tt.want(got)
		})
	}
}

func Test_userManagementServiceLogic_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		credentials model.LoginCredentials
		setup       func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue)
		want        func(*respModel.Response)
	}{
		{
			name: "Success :: Login",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: true})
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				mockDs.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockJwtSvc.EXPECT().GenerateToken(jwt.SigningMethodHS256, "123", time.Nanosecond*6).Return("", nil)
				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: Unable to fetch user from db",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: true})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(nil, errors.New(""))

				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFetchingUser),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: No user found",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: true})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(nil, nil)

				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusUnauthorized,
					Message: codes.GetErr(codes.AccNotFound),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: Get from db failure",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123"})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(users, nil)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(nil, errors.New(""))
				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrFetchingUser),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: Password doesnt match",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123"})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(users, nil)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(nil, nil)
				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusUnauthorized,
					Message: codes.GetErr(codes.InvaliCredentials),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: Acc. Activation in progress",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: false, Salt: "12345cd"})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(users, nil)
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				//	mockJwtSvc.EXPECT().GenerateToken(jwt.SigningMethodHS256, "123", time.Nanosecond*6).Return("", errors.New(""))

				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusAccepted,
					Message: codes.GetErr(codes.ErrLogging),
					Data:    codes.GetErr(codes.AccActivationInProcess),
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: unable to generate token",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: true, Salt: "12345cd"})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(users, nil)
				mockJwtSvc.EXPECT().GenerateToken(jwt.SigningMethodHS256, "123", time.Nanosecond*6).Return("", errors.New(""))
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrGenerateJwt),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
		{
			name: "Failure :: Login :: Update failure",
			credentials: model.LoginCredentials{
				Email:    "vatsal@gmail.com",
				Password: "Abcde@12345",
			},
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockDs := mock.NewMockDataSourceI(mockCtrl)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				var users []model.User
				users = append(users, model.User{Email: "vatsal@gmail.com", Id: "123", Active: true})
				mockDs.EXPECT().Get(map[string]interface{}{"email": "vatsal@gmail.com"}).Times(1).Return(users, nil)
				mockJwtSvc.EXPECT().GenerateToken(jwt.SigningMethodHS256, "123", time.Nanosecond*6).Return("", nil)
				mockDs.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(errors.New(""))
				mockDs.EXPECT().Get(gomock.Any()).Times(1).Return(users, nil)
				return mockDs, mockJwtSvc, config.MsgQueue{MsgBroker: sdk.NewMsgBrokerSvc("http://localhost:9095")}
			},
			want: func(resp *respModel.Response) {
				temp := respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrLogging),
					Data:    nil,
				}
				if !reflect.DeepEqual(resp, &temp) {
					t.Errorf("Want: %v, Got: %v", temp, resp)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewUserMgmtSvcLogic(tt.setup())

			got := rec.Login(httptest.NewRecorder(), tt.credentials)

			tt.want(got)
		})
	}
}
