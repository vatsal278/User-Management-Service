package handler

import (
	"bytes"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/repo/authentication"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
	"net/http/httptest"
	"testing"

	"github.com/PereRohit/util/testutil"
	"github.com/golang/mock/gomock"

	"github.com/vatsal278/UserManagementService/pkg/mock"
)

func Test_userManagementService_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name        string
		setup       func() userMgmtSvc
		wantSvcName string
		wantMsg     string
		wantStat    bool
	}{
		{
			name: "Success",
			setup: func() userMgmtSvc {
				mockLogic := mock.NewMockUserMgmtSvcLogicIer(mockCtrl)
				mockLogic.EXPECT().HealthCheck().
					Return(true).Times(1)

				rec := userMgmtSvc{
					logic: mockLogic,
				}

				return rec
			},
			wantSvcName: UserManagementServiceName,
			wantMsg:     "",
			wantStat:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.setup()

			svcName, msg, stat := receiver.HealthCheck()

			diff := testutil.Diff(svcName, tt.wantSvcName)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(msg, tt.wantMsg)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}

			diff = testutil.Diff(stat, tt.wantStat)
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}

func Test_userManagementServiceLogic_HealthCheck(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name  string
		setup func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue)
		want  func(recorder httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setup: func() (datasource.DataSourceI, authentication.JWTService, config.MsgQueue) {
				mockLogic := mock.NewMockUserManagementServiceHandler(mockCtrl)
				mockLogic.
					mockDs.EXPECT().HealthCheck().Times(1).
					Return(true)

				return mockDs, nil, config.MsgQueue{}
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			rec := NewUserMgmtSvc(tt.setup())
			rec.SignUp(w, r)
			want(w)

		})
	}
}
