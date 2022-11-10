package config

import (
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
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
				},
				ServiceRouteVersion: "v2",
				SvrCfg:              config.ServerConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitSvcConfig(tt.args.cfg)
			diff := testutil.Diff(got, tt.want)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
