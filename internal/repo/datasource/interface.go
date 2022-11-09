package datasource

import (
	"github.com/vatsal278/UserManagementService/internal/model"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_datasource.go --package=mock github.com/vatsal278/user-mgmt-svc/internal/repo/datasource DataSource

type DataSourceI interface {
	HealthCheck() bool
	Get(map[string]interface{}) ([]model.User, error)
	Insert(user model.User) error
	Update(filterSet map[string]interface{}, filterWhere map[string]interface{}) error
}
