package jwtSvc

import (
	"fmt"
	"github.com/PereRohit/util/log"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_jwt.go --package=mock github.com/vatsal278/UserManagementService/internal/repo/jwt JWTService

type JWTService interface {
	GenerateToken(signingMethod jwt.SigningMethod, userId string, validity time.Duration) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type authCustomClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	userId    string
}

func JWTAuthService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "DefaultSecretJwtKey"
	}
	return secret
}

func (service *jwtService) GenerateToken(signingMethod jwt.SigningMethod, userId string, validity time.Duration) (string, error) {
	claims := &authCustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30 * validity).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(signingMethod, claims)

	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return t, nil
}

func (service *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			err := fmt.Errorf("invalid token %+v", token.Header["alg"])
			return nil, err
		}
		return []byte(service.secretKey), nil
	})

}
