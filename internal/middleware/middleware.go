package middleware

import (
	"github.com/PereRohit/util/response"
	"github.com/dgrijalva/jwt-go"
	svcCfg "github.com/vatsal278/UserManagementService/internal/config"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/jwt"
	"github.com/vatsal278/UserManagementService/pkg/session"
	"net/http"
	"strings"
)

type UserMgmtMiddleware struct {
	cfg *svcCfg.SvcConfig
	jwt jwtSvc.JWTService
}

func NewUserMgmtMiddleware(cfg *svcCfg.SvcConfig) *UserMgmtMiddleware {
	return &UserMgmtMiddleware{
		cfg: cfg,
	}
}

func (u UserMgmtMiddleware) ExtractUser(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			response.ToJson(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		if cookie.Value == "" {
			response.ToJson(w, http.StatusUnauthorized, "UnAuthorized", nil)
			return
		}
		token, err := u.cfg.JwtSvc.JwtSvc.ValidateToken(cookie.Value)
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
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.ToJson(w, http.StatusInternalServerError, "unable to assert claims", nil)
			return
		}
		userId := mapClaims["user_id"]
		ctx := session.SetSession(r.Context(), userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u UserMgmtMiddleware) ScreenRequest(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var urlMatch bool
		if u.cfg.Cfg.MessageQueue.UrlCheck != false {
			if r.UserAgent() != u.cfg.Cfg.MessageQueue.UserAgent {
				response.ToJson(w, http.StatusUnauthorized, "UnAuthorized user agent", nil)
				return
			}
			for _, v := range u.cfg.Cfg.MessageQueue.AllowedUrl {
				if v == r.RemoteAddr {
					urlMatch = true
				}
			}
			if urlMatch != true {
				response.ToJson(w, http.StatusUnauthorized, "UnAuthorized url", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
