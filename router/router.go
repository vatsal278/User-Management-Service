package router

import (
	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/middleware"
	"github.com/gorilla/mux"
	"github.com/vatsal278/user-mgmt-svc/internal/config"
	"github.com/vatsal278/user-mgmt-svc/internal/handler"
	middleware2 "github.com/vatsal278/user-mgmt-svc/internal/middleware"
	sqlDb "github.com/vatsal278/user-mgmt-svc/internal/repo/datasource"
	jwtSvc "github.com/vatsal278/user-mgmt-svc/internal/repo/jwt"
	"net/http"
)

func Register(svcCfg *config.SvcConfig) *mux.Router {
	m := mux.NewRouter()

	// group all routes for specific version. e.g.: /v1
	if svcCfg.ServiceRouteVersion != "" {
		m = m.PathPrefix("/" + svcCfg.ServiceRouteVersion).Subrouter()
	}

	m.StrictSlash(true)
	m.Use(middleware.RequestHijacker)
	m.Use(middleware.RecoverPanic)

	commons := handler.NewCommonSvc()
	m.HandleFunc(constant.HealthRoute, commons.HealthCheck).Methods(http.MethodGet)
	m.NotFoundHandler = http.HandlerFunc(commons.RouteNotFound)
	m.MethodNotAllowedHandler = http.HandlerFunc(commons.MethodNotAllowed)

	// attach routes for services below
	m = attachUserMgmtSvcRoutes(m, svcCfg)

	return m
}

func attachUserMgmtSvcRoutes(m *mux.Router, svcCfg *config.SvcConfig) *mux.Router {
	dbSvc := sqlDb.NewSql(svcCfg.DbSvc, svcCfg.DbCfg.TableName)
	jwtService := jwtSvc.JWTAuthService()
	loginService := jwtSvc.StaticLoginService()

	svc := handler.NewUserMgmtSvc(dbSvc, loginService, jwtService)
	middleware := middleware2.NewUserMgmtMiddleware(svcCfg.Cfg)

	m.HandleFunc("/register", svc.SignUp).Methods(http.MethodPost)
	m.HandleFunc("/login", svc.Login).Methods(http.MethodPost)

	route1 := m.PathPrefix("/activate").Subrouter()
	route1.HandleFunc("", svc.Activation).Methods(http.MethodPut)
	route1.Use(middleware.ScreenRequest)

	route2 := m.PathPrefix("/user").Subrouter()
	route2.HandleFunc("", svc.UserDetails).Methods(http.MethodGet)
	route2.Use(middleware.ExtractUser)

	return m
}
