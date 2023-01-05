package config

import (
	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/response"
	"github.com/PereRohit/util/testutil"
	"github.com/gorilla/mux"
	jwtSvc "github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/msgbroker/pkg/crypt"
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
						MessageQueue: MsgQueueCfg{SvcUrl: srv.URL, ActivatedAccountChannel: "account.activation.channel", NewAccountChannel: "new.account.channel", Key: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb1FJQkFBS0NBUUVBejdiaXVpck01TEJjN2pzeTFwS0dpNkZsNm0zWXYrZnFzYzhHMFNIQlFRQlUyYmJHCmI4MHhpM1h0bW1zVURkWGhyNXhaMnA0c0Q4T0J2QjV0eWpZZmJyY1oyc2gzcVp1Rm1EZExpbGJjSUVKdUdaM2UKTWRrZWpRRFdTKzJvMmFyQXI5dFBqTGVKTXk4THhUVVlKNmw4NnFTVVl0aDJNMWtXcUsrc2hWK3hlNnI4NjR0VApGeENvWEp1NU1ILys4Y3ZLeTlIMHh5OHR1V0JKVGF0V2lJd3pqaGlEU1NEclZiMVA2UDJYL2NZOGo5WmJQZW9hCm5aV3VkNlhDYnZJN1Faa1lSMFpreTEyb1grZWxTa2Vpa1pkQWxQbmh6OTNnTUw5cUhYWmIvdC9YWXZOSXRjalMKVzA4cFBXNGUvN2NsMkJVMHBaN2g1dWVEenBjSFhNbUd1WmlPT3dJREFRQUJBb0lCQVFDTHhsRFo0QlZTeXM4dQpUTTNRRUhmVG5EOWR1cCtCdkFsMXI0K3h5Vm9uYUphd2pzc0h6dmZKRmdsV3dUbVVlZG5OOTVPTGhxYTEwT1VMCmR4cUFXVjFiZm9FNmRXMzR4enZtQzBlZEJ3aEgrUXZuMXhEL1VGQzdwOVdNOEplUUtkUlNRbTFNanZFWGJWQXAKVzZvdWZtSWQ3N1Flcy9VT1pxUFZ6YWwxY3NpWEkyanhEd3F4ZkQ1Mk5CUUZqaXpLeXhJUHRxU2xNTGhjZG8rYQpwZzg1RTBDVTNPSGRXS0dack0vbjJ1UGhiRnRNWTNDOTVwTjh6eERHUnpQNlByME9xL1MrTHlid1FydUpNUU1CCmZ2enVPUGNHSzUrcVhlSlBPSzB1ZnFhbUhEQ2tHanJHZGRHZmZ0M0lMbzRCV0xrM1I4TXBVYUdoR0hTeng3QVoKNmtCOHhkUTVBb0dCQVBVNTlmOWNTcHh3ZXd4clY4Nm13OWRMTmNZMEsvVXhsME9EQlFBMlVleGJETFo4VnZocwpCcXRzZjdQZFRpdU53WEc0ZnE0bU5SVnJ1K2NyYWhlREI3ZTBqdk9LOGdpQ0QxdlhsZzN5WDZhbDZsUTQ1OXJYCkpHVmExc2hkRjBOdkFYTXR0N2FzQW1QK1kzRXlDSnpQK2NZOFZGaExJRHdlMFJzYzdJQ29Ma1Z2QW9HQkFOalgKQjlGVS8xYy8yUk9UaDdZYldjdHh2dDBtaDMrcStZaU9jQnZsQUMzbXNxc0NNTSt1ZGZ3U2FNQlpoSXhxUWRYUgpTQTF0VFgyY1hKZ0VpZ1RGNnJBNGdudmxOZW5QTGJvNlJydmZoNkdsejdIZXZ4cGY1ZWtYNHRTM3FEcWFWSDBGClZvY25jSzc3L21DWmM0VVVHeDJDTHYvbFNtZzRzaUtCQTltZFFoWDFBbjlyU1BCV3lBbmNaMWx1RlloVTRLRE4Ka0JuMm5OeWVhUlBFZFkyNmlnbE5Yb2d4VGpTK2VvUndld2RqcVc2Sm4zc0NSYlVtZTVDOXptUm12cGVyc2FldQp0MC9UUFBhbXdqLzE3bHUzdmxJYWxudnVYUGNTeHcwbFNwaXRFQTBkYzNNdThORnZHZEh4N1ZtVUxFK1lTMlQ3ClZXbVJOMHpqQUpoN1JDdzBIV0FoQW9HQUNHQ21hS3dFQVhieUNCT1hGcTRQMWhCYTgyaGRxODBMUHY5aHpYSVgKZzY1NkVLbFJBWFVZRWRrVU92bzZhTUppTU1TWktBdWxCc2xYdW5mU2JVVElRRzZ1ZStMckpsRmV6dWNaZklDeQpXTWh6TWNnTlVoT0thbXNGMUhvVUFjK2NuQWZzdytQK01vU0IyM0dTU1AzeDNqMzlXdDJjOWxIYWNBTFVCMEJRCklWRUNnWUJETHhjM1o3Y1o3VU5PNndRdy8xTWRGZXEwRGJ6dlY0Z1hlVkpoUHZWeU1RR2hiZjVESTFXMngwWGgKUjBHeUVyZTIrVGVXZEswRUZCS3FwTVRlSFNZNmM2QWY1UUZmOTJsRW1PSlVHOVhXQ3FBS3pqMHFUY0l1bldsWgpsaXpWdjJIM0hDbmlTWXMzbG9kSFFsSzcyZTBzcC9Kd25QN2hjMHlsNmh1cHdDeVQ3UT09Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t"},
						Cookie:       CookieStruct{ExpiryStr: "5m"},
					},
				}
			},
			want: func(arg args) *SvcConfig {
				s := sdk.NewMsgBrokerSvc(arg.cfg.MessageQueue.SvcUrl)
				privateKey, _ := crypt.PEMStrAsPrivKey(arg.cfg.MessageQueue.Key)
				required := &SvcConfig{
					JwtSvc: JWTSvc{JwtSvc: jwtSvc.JWTAuthService("")},
					Cfg: &Config{
						ServiceRouteVersion: "v2",
						ServerConfig:        config.ServerConfig{},
						DataBase: DbCfg{
							Driver: "mysql",
						},
						MessageQueue: MsgQueueCfg{SvcUrl: arg.cfg.MessageQueue.SvcUrl, NewAccountChannel: "new.account.channel", ActivatedAccountChannel: "account.activation.channel", Key: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb1FJQkFBS0NBUUVBejdiaXVpck01TEJjN2pzeTFwS0dpNkZsNm0zWXYrZnFzYzhHMFNIQlFRQlUyYmJHCmI4MHhpM1h0bW1zVURkWGhyNXhaMnA0c0Q4T0J2QjV0eWpZZmJyY1oyc2gzcVp1Rm1EZExpbGJjSUVKdUdaM2UKTWRrZWpRRFdTKzJvMmFyQXI5dFBqTGVKTXk4THhUVVlKNmw4NnFTVVl0aDJNMWtXcUsrc2hWK3hlNnI4NjR0VApGeENvWEp1NU1ILys4Y3ZLeTlIMHh5OHR1V0JKVGF0V2lJd3pqaGlEU1NEclZiMVA2UDJYL2NZOGo5WmJQZW9hCm5aV3VkNlhDYnZJN1Faa1lSMFpreTEyb1grZWxTa2Vpa1pkQWxQbmh6OTNnTUw5cUhYWmIvdC9YWXZOSXRjalMKVzA4cFBXNGUvN2NsMkJVMHBaN2g1dWVEenBjSFhNbUd1WmlPT3dJREFRQUJBb0lCQVFDTHhsRFo0QlZTeXM4dQpUTTNRRUhmVG5EOWR1cCtCdkFsMXI0K3h5Vm9uYUphd2pzc0h6dmZKRmdsV3dUbVVlZG5OOTVPTGhxYTEwT1VMCmR4cUFXVjFiZm9FNmRXMzR4enZtQzBlZEJ3aEgrUXZuMXhEL1VGQzdwOVdNOEplUUtkUlNRbTFNanZFWGJWQXAKVzZvdWZtSWQ3N1Flcy9VT1pxUFZ6YWwxY3NpWEkyanhEd3F4ZkQ1Mk5CUUZqaXpLeXhJUHRxU2xNTGhjZG8rYQpwZzg1RTBDVTNPSGRXS0dack0vbjJ1UGhiRnRNWTNDOTVwTjh6eERHUnpQNlByME9xL1MrTHlid1FydUpNUU1CCmZ2enVPUGNHSzUrcVhlSlBPSzB1ZnFhbUhEQ2tHanJHZGRHZmZ0M0lMbzRCV0xrM1I4TXBVYUdoR0hTeng3QVoKNmtCOHhkUTVBb0dCQVBVNTlmOWNTcHh3ZXd4clY4Nm13OWRMTmNZMEsvVXhsME9EQlFBMlVleGJETFo4VnZocwpCcXRzZjdQZFRpdU53WEc0ZnE0bU5SVnJ1K2NyYWhlREI3ZTBqdk9LOGdpQ0QxdlhsZzN5WDZhbDZsUTQ1OXJYCkpHVmExc2hkRjBOdkFYTXR0N2FzQW1QK1kzRXlDSnpQK2NZOFZGaExJRHdlMFJzYzdJQ29Ma1Z2QW9HQkFOalgKQjlGVS8xYy8yUk9UaDdZYldjdHh2dDBtaDMrcStZaU9jQnZsQUMzbXNxc0NNTSt1ZGZ3U2FNQlpoSXhxUWRYUgpTQTF0VFgyY1hKZ0VpZ1RGNnJBNGdudmxOZW5QTGJvNlJydmZoNkdsejdIZXZ4cGY1ZWtYNHRTM3FEcWFWSDBGClZvY25jSzc3L21DWmM0VVVHeDJDTHYvbFNtZzRzaUtCQTltZFFoWDFBbjlyU1BCV3lBbmNaMWx1RlloVTRLRE4Ka0JuMm5OeWVhUlBFZFkyNmlnbE5Yb2d4VGpTK2VvUndld2RqcVc2Sm4zc0NSYlVtZTVDOXptUm12cGVyc2FldQp0MC9UUFBhbXdqLzE3bHUzdmxJYWxudnVYUGNTeHcwbFNwaXRFQTBkYzNNdThORnZHZEh4N1ZtVUxFK1lTMlQ3ClZXbVJOMHpqQUpoN1JDdzBIV0FoQW9HQUNHQ21hS3dFQVhieUNCT1hGcTRQMWhCYTgyaGRxODBMUHY5aHpYSVgKZzY1NkVLbFJBWFVZRWRrVU92bzZhTUppTU1TWktBdWxCc2xYdW5mU2JVVElRRzZ1ZStMckpsRmV6dWNaZklDeQpXTWh6TWNnTlVoT0thbXNGMUhvVUFjK2NuQWZzdytQK01vU0IyM0dTU1AzeDNqMzlXdDJjOWxIYWNBTFVCMEJRCklWRUNnWUJETHhjM1o3Y1o3VU5PNndRdy8xTWRGZXEwRGJ6dlY0Z1hlVkpoUHZWeU1RR2hiZjVESTFXMngwWGgKUjBHeUVyZTIrVGVXZEswRUZCS3FwTVRlSFNZNmM2QWY1UUZmOTJsRW1PSlVHOVhXQ3FBS3pqMHFUY0l1bldsWgpsaXpWdjJIM0hDbmlTWXMzbG9kSFFsSzcyZTBzcC9Kd25QN2hjMHlsNmh1cHdDeVQ3UT09Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t"},
						Cookie:       CookieStruct{Expiry: time.Minute * 5, ExpiryStr: "5m"},
					},
					ServiceRouteVersion: "v2",
					SvrCfg:              config.ServerConfig{},
					MsgBrokerSvc:        MsgQueue{MsgBroker: s, PubId: "123", Channel: "new.account.channel", PrivateKey: *privateKey},
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
