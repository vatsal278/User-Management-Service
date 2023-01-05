package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	respModel "github.com/PereRohit/util/model"
	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"
	"github.com/vatsal278/UserManagementService/internal/codes"
	"github.com/vatsal278/UserManagementService/internal/model"
	"github.com/vatsal278/UserManagementService/pkg/session"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/UserManagementService/pkg/mock"
)

type Reader string

func (Reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test_userManagementService_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		setup       func() userMgmtSvc
		wantSvcName string
		wantMsg     string
		wantStat    bool
	}{
		{
			name: "Success",
			setup: func() userMgmtSvc {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().HealthCheck().
					Return(true).Times(1)

				rec := userMgmtSvc{
					logic: mockLogic,
				}

				return rec
			},
			wantSvcName: UserManagementServiceName,
			wantMsg:     "",
			wantStat:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.setup()

			svcName, msg, stat := receiver.HealthCheck()

			diff := testutil.Diff(svcName, tt.wantSvcName)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(msg, tt.wantMsg)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(stat, tt.wantStat)
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
		name  string
		model model.SignUpCredentials
		setup func() (*userMgmtSvc, *http.Request)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			model: model.SignUpCredentials{
				Name:             "Vatsal",
				Email:            "vatsal@gmail.com",
				Password:         "Abc@123",
				RegistrationDate: "15-11-2022 00:00:00",
			},
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().Signup(gomock.Any()).Times(1).Return(&respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    codes.GetErr(codes.AccActivationInProcess),
				})
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.SignUpCredentials{
					Name:             "Vatsal",
					Email:            "vatsal@gmail.com",
					Password:         "AbcDe@123",
					RegistrationDate: "15-11-2022 00:00:00",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusCreated,
					Message: codes.GetErr(codes.Success),
					Data:    codes.GetErr(codes.AccActivationInProcess),
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}

			},
		},
		{
			name: "Failure :: signUp:: Read all failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/register", Reader(""))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: "request body read : test error",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: signUp:: json unmarshall failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("")))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "put data into data: unexpected end of JSON input",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: SignUp:: Validate failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.SignUpCredentials{
					Name:             "",
					Email:            "vatsal@gmail.com",
					Password:         "Abc@123",
					RegistrationDate: "15-11-2022 00:00:00",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "validation 1: field <Name> with value <> failed for <required> validation.\n2: field <Password> with value <Abc@123> failed for <min> validation.\n",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: SignUp:: Password Validate:: numeric",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.SignUpCredentials{
					Name:             "vatsal",
					Email:            "vatsal@gmail.com",
					Password:         "Abcdefghj",
					RegistrationDate: "15-11-2022 00:00:00",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrPassNumeric),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: SignUp:: Password Validate :: special char",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.SignUpCredentials{
					Name:             "vatsal",
					Email:            "vatsal@gmail.com",
					Password:         "Abcdefghj1",
					RegistrationDate: "15-11-2022 00:00:00",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrPassSpecial),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: SignUp :: time parse failure",
			model: model.SignUpCredentials{
				Name:             "Vatsal",
				Email:            "vatsal@gmail.com",
				Password:         "Abc@123",
				RegistrationDate: "15-11-2022 00:00:00",
			},
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.SignUpCredentials{
					Name:             "Vatsal",
					Email:            "vatsal@gmail.com",
					Password:         "Abcde@123",
					RegistrationDate: "ABC",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrParseRegDate),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.SignUp(w, r)
			tt.want(*w)
		})
	}
}

func Test_userManagementServiceLogic_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		model model.LoginCredentials
		setup func() (*userMgmtSvc, *http.Request)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return(&respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				})
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.LoginCredentials{
					Email:    "vatsal@gmail.com",
					Password: "Abcde@123",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/login", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: Read all failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/login", Reader(""))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: "request body read : test error",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: json unmarshall failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("")))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "put data into data: unexpected end of JSON input",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: Validate failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.LoginCredentials{
					Email:    "Vatsal@gmail.com",
					Password: "A@123",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "validation 1: field <Password> with value <A@123> failed for <min> validation.\n",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: Password Validate failure:: upper case",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.LoginCredentials{
					Email:    "Vatsal@gmail.com",
					Password: "aasdfghjkl",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrPassUpperCase),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: Password Validate failure:: lowercase",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(model.LoginCredentials{
					Email:    "Vatsal@gmail.com",
					Password: "ASDFGHJKL",
				})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrPassLowerCase),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.Login(w, r)
			tt.want(*w)
		})
	}
}
func Test_userManagementServiceLogic_Activate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		model model.LoginCredentials
		setup func() (*userMgmtSvc, *http.Request)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().Activate(gomock.Any()).Times(1).Return(&respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				})
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				by, err := json.Marshal(map[string]interface{}{"user_id": "123"})
				if err != nil {
					t.Fail()
				}
				r := httptest.NewRequest("PUT", "/activate", bytes.NewBuffer(by))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Activate:: Read all failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("PUT", "/activate", Reader(""))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: "request body read : test error",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure :: Login:: json unmarshall failure",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("PUT", "/activate", bytes.NewBuffer([]byte("")))
				return svc, r
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: "put data into data: unexpected end of JSON input",
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.Activation(w, r)
			tt.want(*w)
		})
	}
}
func Test_userManagementServiceLogic_User(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		model model.LoginCredentials
		setup func() (*userMgmtSvc, *http.Request)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().UserData(gomock.Any()).Times(1).Return(&respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				})
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("PUT", "/activate", nil)
				ctx := session.SetSession(r.Context(), "1234")
				return svc, r.WithContext(ctx)
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusOK,
					Message: codes.GetErr(codes.Success),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure:: logic :: internal server error",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().UserData("1234").Return(&respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrAssertUserid),
					Data:    nil,
				})
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("PUT", "/activate", nil)
				ctx := session.SetSession(r.Context(), "1234")
				return svc, r.WithContext(ctx)
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusInternalServerError,
					Message: codes.GetErr(codes.ErrAssertUserid),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
		{
			name: "Failure:: err asserting to string",
			setup: func() (*userMgmtSvc, *http.Request) {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				svc := &userMgmtSvc{
					logic: mockLogic,
				}
				r := httptest.NewRequest("PUT", "/activate", nil)
				ctx := session.SetSession(r.Context(), 1.11)
				return svc, r.WithContext(ctx)
			},
			want: func(rec httptest.ResponseRecorder) {
				b, err := ioutil.ReadAll(rec.Body)
				if err != nil {
					return
				}
				var response respModel.Response
				err = json.Unmarshal(b, &response)
				tempResp := &respModel.Response{
					Status:  http.StatusBadRequest,
					Message: codes.GetErr(codes.ErrAssertUserid),
					Data:    nil,
				}
				if !reflect.DeepEqual(&response, tempResp) {
					t.Errorf("Want: %v, Got: %v", tempResp, &response)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			x, r := tt.setup()
			x.UserDetails(w, r)
			tt.want(*w)
		})
	}
}
