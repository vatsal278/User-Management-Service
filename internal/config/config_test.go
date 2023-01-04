package config

import (
	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/response"
	"github.com/PereRohit/util/testutil"
	"github.com/gorilla/mux"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/msgbroker/pkg/sdk"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInitSvcConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args func() args
		want func(args) *SvcConfig
	}{
		{
			name: "Success",
			args: func() args {
				router := mux.NewRouter()
				router.HandleFunc("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusOK, "Success", map[string]interface{}{"id": "123"})
				})
				router.HandleFunc("/register/subscriber", func(w http.ResponseWriter, r *http.Request) {
					response.ToJson(w, http.StatusOK, "Success", nil)
				})
				srv := httptest.NewServer(router)
				return args{
					cfg: Config{
						ServiceRouteVersion: "v2",
						ServerConfig:        config.ServerConfig{},
						DataBase: DbCfg{
							Driver: "mysql",
						},
						MessageQueue: MsgQueueCfg{SvcUrl: srv.URL, ActivatedAccountChannel: "account.activation.channel", NewAccountChannel: "new.account.channel"},
						Cookie:       CookieStruct{ExpiryStr: "5m"},
					},
				}
			},
			want: func(arg args) *SvcConfig {
				s := sdk.NewMsgBrokerSvc(arg.cfg.MessageQueue.SvcUrl)
				required := &SvcConfig{
					JwtSvc: JWTSvc{JwtSvc: jwtSvc.JWTAuthService("")},
					Cfg: &Config{
						ServiceRouteVersion: "v2",
						ServerConfig:        config.ServerConfig{},
						DataBase: DbCfg{
							Driver: "mysql",
						},
						MessageQueue: MsgQueueCfg{SvcUrl: arg.cfg.MessageQueue.SvcUrl, NewAccountChannel: "new.account.channel", ActivatedAccountChannel: "account.activation.channel"},
						Cookie:       CookieStruct{Expiry: time.Minute * 5, ExpiryStr: "5m"},
					},
					ServiceRouteVersion: "v2",
					SvrCfg:              config.ServerConfig{},
					MsgBrokerSvc:        MsgQueue{MsgBroker: s, PubId: "123", Channel: "new.account.channel"},
				}
				return required
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.args()
			got := InitSvcConfig(s.cfg)

			diff := testutil.Diff(got, tt.want(s))
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
