package middleware

import (
	"encoding/json"
	"errors"
	"github.com/PereRohit/util/model"
	"github.com/PereRohit/util/response"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/pkg/mock"
	"github.com/vatsal278/UserManagementService/pkg/session"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var x func()

var hit = false

func test(w http.ResponseWriter, r *http.Request) {
	hit = true
	c := r.Context()
	id := session.GetSession(c)
	response.ToJson(w, http.StatusOK, "passed", id)
}

func TestUserMgmtMiddleware_ExtractUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name      string
		setupFunc func() (*http.Request, authentication.JWTService)
		validator func(*httptest.ResponseRecorder)
	}{
		{
			name: "SUCCESS::ExtractUser",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)

				token := jwtGo.Token{
					Claims: jwtGo.MapClaims{"user_id": "123"},
					Valid:  true,
				}
				mockJwtSvc.EXPECT().ValidateToken("jwtToken").Return(&token, nil)
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "jwtToken",
				})
				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != true {
					t.Errorf("Want: %v, Got: %v", true, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					return
				}
				expected := &model.Response{
					Status:  http.StatusOK,
					Message: "passed",
					Data:    "123",
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: empty token value",
			setupFunc: func() (*http.Request, authentication.JWTService) {
				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "",
				})
				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "UnAuthorized",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: no cookie found",
			setupFunc: func() (*http.Request, authentication.JWTService) {
				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "UnAuthorized",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: compared literals not same",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: " jwtToken",
				})
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				err := errors.New(" err ")
				mockJwtSvc.EXPECT().ValidateToken(gomock.Any()).Return(nil, err)

				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "Compared literals are not same",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: Token is expired ",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				//authentication := jwtSvc.JWTAuthService()
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "123",
				})
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				err := errors.New(" Token is expired ")
				mockJwtSvc.EXPECT().ValidateToken(gomock.Any()).Return(nil, err)

				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "Token is expired",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: Token is invalid ",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "123",
				})
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				token := jwtGo.Token{Valid: false}
				mockJwtSvc.EXPECT().ValidateToken(gomock.Any()).Return(&token, nil)

				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: mapClaims not ok",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				//authentication := jwtSvc.JWTAuthService()
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "123",
				})
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				token := jwtGo.Token{
					Claims: nil,
					Valid:  true,
				}
				mockJwtSvc.EXPECT().ValidateToken(gomock.Any()).Return(&token, nil)

				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusInternalServerError,
					Message: "unable to assert claims",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::ExtractUser:: user id not in claims",
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "123",
				})
				mockJwtSvc := mock.NewMockJWTService(mockCtrl)
				token := jwtGo.Token{
					Claims: jwtGo.MapClaims{},
					Valid:  true,
				}
				mockJwtSvc.EXPECT().ValidateToken(gomock.Any()).Return(&token, nil)

				return req, mockJwtSvc
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusInternalServerError,
					Message: "unable to assert userid",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: seting up all instances for the specific test case
			res := httptest.NewRecorder()
			req, jwt := tt.setupFunc()

			// STEP 2: call the test function
			middleware := NewUserMgmtMiddleware(&config.SvcConfig{
				JwtSvc: config.JWTSvc{
					JwtSvc: jwt,
				},
			})
			hit = false
			x := middleware.ExtractUser(test)
			x.ServeHTTP(res, req)

			tt.validator(res)

		})
	}
}
func TestUserMgmtMiddleware_ScreenRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name      string
		config    config.Config
		setupFunc func() (*http.Request, authentication.JWTService)
		validator func(*httptest.ResponseRecorder)
	}{
		{
			name: "SUCCESS::Screen Request",
			config: config.Config{
				MessageQueue: config.MsgQueue{
					AllowedUrl: []string{"192.0.2.1:1234", "value2", "value3"},
					UserAgent:  "UserAgent",
					UrlCheck:   true,
				}},
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				jwt := mock.NewMockJWTService(mockCtrl)
				req.Header.Set("User-Agent", "UserAgent")
				return req, jwt
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != true {
					t.Errorf("Want: %v, Got: %v", true, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusOK,
					Message: "passed",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::Screen Request:: unauthorized user agent",
			config: config.Config{
				MessageQueue: config.MsgQueue{
					AllowedUrl: []string{"192.0.2.1:1234", "value2", "value3"},
					UserAgent:  "U",
					UrlCheck:   true,
				}},
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				jwt := mock.NewMockJWTService(mockCtrl)
				req.Header.Set("User-Agent", "UserAgent")
				return req, jwt
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "UnAuthorized user agent",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Failure::Screen Request:: unauthorized url",
			config: config.Config{
				MessageQueue: config.MsgQueue{
					AllowedUrl: []string{"value", "value2", "value3"},
					UserAgent:  "UserAgent",
					UrlCheck:   true,
				}},
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				jwt := mock.NewMockJWTService(mockCtrl)
				req.Header.Set("User-Agent", "UserAgent")
				return req, jwt
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != false {
					t.Errorf("Want: %v, Got: %v", false, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusUnauthorized,
					Message: "UnAuthorized url",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
		{
			name: "Success::Screen Request:: url check not required",
			config: config.Config{
				MessageQueue: config.MsgQueue{
					UrlCheck:  false,
					UserAgent: "UserAgent",
				}},
			setupFunc: func() (*http.Request, authentication.JWTService) {

				req := httptest.NewRequest(http.MethodGet, "http://localhost:80", nil)
				jwt := mock.NewMockJWTService(mockCtrl)
				req.Header.Set("User-Agent", "UserAgent")
				return req, jwt
			},
			validator: func(res *httptest.ResponseRecorder) {
				if hit != true {
					t.Errorf("Want: %v, Got: %v", true, hit)
				}
				by, _ := ioutil.ReadAll(res.Body)
				result := model.Response{}
				err := json.Unmarshal(by, &result)
				if err != nil {
					t.Log(err)
					t.Fail()
					return
				}
				expected := &model.Response{
					Status:  http.StatusOK,
					Message: "passed",
					Data:    nil,
				}
				if !reflect.DeepEqual(&result, expected) {
					t.Errorf("Want: %v, Got: %v", expected, result)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// STEP 1: seting up all instances for the specific test case
			res := httptest.NewRecorder()
			req, jwt := tt.setupFunc()

			middleware := NewUserMgmtMiddleware(&config.SvcConfig{
				Cfg: &tt.config,
				JwtSvc: config.JWTSvc{
					JwtSvc: jwt,
				},
			})
			hit = false
			x := middleware.ScreenRequest(test)
			x.ServeHTTP(res, req)

			tt.validator(res)

		})
	}
}
