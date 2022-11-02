package middleware

import (
	"context"
	"errors"
	"github.com/PereRohit/util/response"
	"github.com/dgrijalva/jwt-go"
	svcCfg "github.com/vatsal278/UserManagementService/internal/config"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/jwt"

	"net/http"
	"strings"
)

type UserMgmtMiddleware struct {
	cfg *svcCfg.Config
}

func NewUserMgmtMiddleware(cfg *svcCfg.Config) *UserMgmtMiddleware {
	return &UserMgmtMiddleware{
		cfg: cfg,
	}
}

type UserId struct{}

func generateContext(r *http.Request, token *jwt.Token) (context.Context, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to assert claims")
	}
	userId := mapClaims["user_id"]
	return context.WithValue(r.Context(), UserId{}, userId), nil
}
func ExtractId(r *http.Request) (string, error) {
	id := r.Context().Value(UserId{})
	i, ok := id.(string)
	if !ok {
		return "", errors.New("cannot assert id")
	}
	return i, nil
}
func (u UserMgmtMiddleware) ExtractUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			return
		}
		if cookie.Value == "" {
			response.ToJson(w, http.StatusUnauthorized, "UnAuthorized", nil)
			return
		}

		token, err := jwtSvc.JWTAuthService().ValidateToken(cookie.Value)
		if err != nil {
			if strings.Contains(err.Error(), "Token is expired") {
				response.ToJson(w, http.StatusUnauthorized, "Token is expired", nil)
				return
			}
			response.ToJson(w, http.StatusUnauthorized, "Compared literals are not same", nil)
			return
		}

		if !token.Valid {
			response.ToJson(w, http.StatusUnauthorized, "Unauthorized", nil)
			return

		}
		ctx, err := generateContext(r, token)
		if err != nil {
			http.Error(w, "Unable to generate context", http.StatusInternalServerError)
			response.ToJson(w, http.StatusInternalServerError, "Unable to generate context", nil)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u UserMgmtMiddleware) ScreenRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u.cfg.MessageQueue != r.RemoteAddr && r.UserAgent() != "message queue" {
			response.ToJson(w, http.StatusInternalServerError, "UnAuthorized user agent", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
