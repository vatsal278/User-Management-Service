package logic

import (
	"net/http"

	"github.com/PereRohit/util/log"
	respModel "github.com/PereRohit/util/model"

	"github.com/vatsal278/UserManagementService/internal/model"
	"github.com/vatsal278/UserManagementService/internal/repo/datasource"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../pkg/mock/mock_logic.go --package=mock github.com/vatsal278/UserManagementService/internal/logic UserManagementServiceLogicIer

type UserManagementServiceLogicIer interface {
	Ping(*model.PingRequest) *respModel.Response
	HealthCheck() bool
}

type userManagementServiceLogic struct {
	dummyDsSvc datasource.DataSource
}

func NewUserManagementServiceLogic(ds datasource.DataSource) UserManagementServiceLogicIer {
	return &userManagementServiceLogic{
		dummyDsSvc: ds,
	}
}

func (l userManagementServiceLogic) Ping(req *model.PingRequest) *respModel.Response {
	// add business logic here
	res, err := l.dummyDsSvc.Ping(&model.PingDs{
		Data: req.Data,
	})
	if err != nil {
		log.Error("datasource error", err)
		return &respModel.Response{
			Status:  http.StatusInternalServerError,
			Message: "",
			Data:    nil,
		}
	}
	return &respModel.Response{
		Status:  http.StatusOK,
		Message: "Pong",
		Data:    res,
	}
}

func (l userManagementServiceLogic) HealthCheck() bool {
	// check all internal services are working fine
	return l.dummyDsSvc.HealthCheck()
}
