package authentication

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestJwtService_GenerateToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name          string
		signingMethod jwt.SigningMethod
		validator     func(JWTService, string, error)
	}{
		{
			name:          "SUCCESS:: Generate Token",
			signingMethod: jwt.SigningMethodHS256,
			validator: func(jwtSvc JWTService, token string, err error) {
				if !reflect.DeepEqual(err, nil) {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
				t.Log(token)
				validateToken, _ := jwtSvc.ValidateToken(token)
				mapClaims, ok := validateToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Log("failed to assert claims")
					t.Fail()
					return
				}
				userId := mapClaims["user_id"]
				t.Log(validateToken.Claims.Valid())
				if !reflect.DeepEqual(userId, "1") {
					t.Errorf("Want: %v, Got: %v", "1", userId)
				}
			},
		},
		{
			name:          "Failure:: Generate Token:: invalid signing method",
			signingMethod: jwt.SigningMethodES256,
			validator: func(jwtSvc JWTService, token string, err error) {
				if !reflect.DeepEqual(err.Error(), errors.New("key is of invalid type").Error()) {
					t.Errorf("Want: %v, Got: %v", "key is of invalid type", err)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtSvc := JWTAuthService("")
			token, err := jwtSvc.GenerateToken(tt.signingMethod, "1", 10)

			tt.validator(jwtSvc, token, err)

		})
	}
}
func TestJwtService_ValidateToken(t *testing.T) {
	jwtSvc := JWTAuthService("")
	tests := []struct {
		name      string
		setupFunc func() string
		validator func(*jwt.Token, error)
	}{
		{
			name: "SUCCESS:: Validate Token",
			setupFunc: func() string {
				token, err := jwtSvc.GenerateToken(jwt.SigningMethodHS256, "1", 360)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				return token
			},
			validator: func(token *jwt.Token, err error) {
				if !reflect.DeepEqual(err, nil) {
					t.Errorf("Want: %v, Got: %v", nil, err)
				}
				mapClaims := token.Claims.(jwt.MapClaims)
				userId := mapClaims["user_id"]
				if !reflect.DeepEqual(userId, "1") {
					t.Errorf("Want: %v, Got: %v", "1", userId)
				}
			},
		},
		{
			name: "Failure:: Validate Token",
			setupFunc: func() string {
				return ""
			},
			validator: func(token *jwt.Token, err error) {
				if !reflect.DeepEqual(err.Error(), errors.New("token contains an invalid number of segments").Error()) {
					t.Errorf("Want: %v, Got: %v", "token contains an invalid number of segments", err)
				}
			},
		},
	}

	// to execute the tests in the table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtSvc.ValidateToken(tt.setupFunc())
			tt.validator(token, err)
		})
	}
}
