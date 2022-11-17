package config

import (
	"github.com/PereRohit/util/response"
	"github.com/gorilla/mux"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/testutil"
)

func TestInitSvcConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
		want *SvcConfig
	}{
		{
			name: "Success",
			args: args{
				cfg: Config{
					ServiceRouteVersion: "v2",
					ServerConfig:        config.ServerConfig{},
					DataBase: DbCfg{
						Driver: "mysql",
					},
					MessageQueue: MsgQueueCfg{SvcUrl: "http://localhost:9091"},
				},
			},
			want: &SvcConfig{
				JwtSvc: JWTSvc{JwtSvc: jwtSvc.JWTAuthService("")},
				Cfg: &Config{
					ServiceRouteVersion: "v2",
					ServerConfig:        config.ServerConfig{},
					DataBase: DbCfg{
						Driver: "mysql",
					},
					MessageQueue: MsgQueueCfg{SvcUrl: "http://localhost:9091"},
				},
				ServiceRouteVersion: "v2",
				SvrCfg:              config.ServerConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
				response.ToJson(w, http.StatusOK, "Success", map[string]interface{}{"id": "123"})
			})
			srv := httptest.NewServer(router)
			defer srv.Close()
			srv.URL = "http://localhost:9091"
			got := InitSvcConfig(tt.args.cfg)
			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
