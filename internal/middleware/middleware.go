package middleware

import (
	"bytes"
	"github.com/PereRohit/util/log"
	"github.com/PereRohit/util/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/vatsal278/UserManagementService/internal/codes"
	svcCfg "github.com/vatsal278/UserManagementService/internal/config"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/pkg/session"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"io"
	"net/http"
	"strings"
)

type UserMgmtMiddleware struct {
	cfg *svcCfg.Config
	jwt jwtSvc.JWTService
	msg func(io.ReadCloser) (string, error)
}

func NewUserMgmtMiddleware(cfg *svcCfg.SvcConfig) *UserMgmtMiddleware {
	msgQueue := sdk.NewMsgBrokerSvc(cfg.Cfg.MessageQueue.SvcUrl)
	msg := msgQueue.ExtractMsg(&cfg.MsgBrokerSvc.PrivateKey)
	return &UserMgmtMiddleware{
		cfg: cfg.Cfg,
		jwt: cfg.JwtSvc.JwtSvc,
		msg: msg,
	}
}

func (u UserMgmtMiddleware) ExtractUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			log.Error(err)
			//user need t login err msg
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		if cookie.Value == "" {
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		token, err := u.jwt.ValidateToken(cookie.Value)
		if err != nil {
			if strings.Contains(err.Error(), "Token is expired") {
				response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrTokenExpired), nil)
				return
			}
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrMatchingToken), nil)
			return
		}
		if !token.Valid {
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorized), nil)
			return
		}
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.ToJson(w, http.StatusInternalServerError, codes.GetErr(codes.ErrAssertClaims), nil)
			return
		}
		userId, ok := mapClaims["user_id"]
		if !ok {
			response.ToJson(w, http.StatusInternalServerError, codes.GetErr(codes.ErrAssertUserid), nil)
			return
		}
		ctx := session.SetSession(r.Context(), userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u UserMgmtMiddleware) ScreenRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var urlMatch bool
		if r.UserAgent() != u.cfg.MessageQueue.UserAgent {
			log.Info(codes.GetErr(codes.ErrUnauthorizedAgent))
			response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorizedAgent), nil)
			return
		}
		if u.cfg.MessageQueue.UrlCheck != false {
			for _, v := range u.cfg.MessageQueue.AllowedUrl {
				if v == r.RemoteAddr {
					urlMatch = true
				}
			}
			if urlMatch != true {
				log.Info(codes.GetErr(codes.ErrUnauthorizedUrl))
				response.ToJson(w, http.StatusUnauthorized, codes.GetErr(codes.ErrUnauthorizedUrl), nil)
				return
			}
		}
		decryptMsg, err := u.msg(r.Body)
		if err != nil {
			log.Error(err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer([]byte(decryptMsg)))
		next.ServeHTTP(w, r)
	})
}
