package datasource

import (
	"database/sql"
	"fmt"
	"github.com/vatsal278/UserManagementService/internal/config"
	"github.com/vatsal278/UserManagementService/internal/model"
	"strings"
)

type sqlDs struct {
	sqlSvc *sql.DB
	table  string
}

func NewSql(dbSvc config.DbSvc, tableName string) DataSourceI {
	return &sqlDs{
		sqlSvc: dbSvc.Db,
		table:  tableName,
	}
}

func (d sqlDs) HealthCheck() bool {
	err := d.sqlSvc.Ping()
	if err != nil {
		return false
	}
	return true
}

func (d sqlDs) Get(filter map[string]interface{}) ([]model.User, error) {
	//order the queries based on email address
	var user model.User
	var users []model.User
	q := fmt.Sprintf("SELECT user_id, email, company_name, name, registered_on, updated_on,salt, active, active_devices FROM %s", d.table)

	filterClause := []string{}

	for k, v := range filter {
		switch v.(type) {
		case string:
			filterClause = append(filterClause, fmt.Sprintf("%s = '%s'", k, v))
		default:
			filterClause = append(filterClause, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClause) > 0 {
		q += fmt.Sprintf(" WHERE %s", strings.Join(filterClause, " AND "))
	}

	q += " ORDER BY email;"
	rows, err := d.sqlSvc.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Email, &user.Company, &user.Name, &user.RegisteredOn, &user.UpdatedOn, &user.Salt, &user.Active, &user.ActiveDevices)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d sqlDs) Insert(user model.User) error {
	queryString := fmt.Sprintf("INSERT INTO %s", d.table)
	_, err := d.sqlSvc.Exec(queryString+"(user_id, name, registered_on, email, salt, company_name, password) VALUES(?, ?,?,?,?, ?,?)", user.Id, user.Name, user.RegisteredOn, user.Email, user.Salt, user.Company, user.Password)
	if err != nil {
		return err
	}
	return err
}

func (d sqlDs) Update(filterSet map[string]interface{}, filterWhere map[string]interface{}) error {
	queryString := fmt.Sprintf("UPDATE %s ", d.table)
	filterClause := []string{}

	for k, v := range filterSet {
		switch v.(type) {
		case string:
			filterClause = append(filterClause, fmt.Sprintf("%s = '%+v'", k, v))
		default:
			filterClause = append(filterClause, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClause) > 0 {
		queryString += fmt.Sprintf(" SET %s", strings.Join(filterClause, " , "))
	}
	filterClauseWhere := []string{}

	for k, v := range filterWhere {
		switch v.(type) {
		case string:
			filterClauseWhere = append(filterClauseWhere, fmt.Sprintf("%s = '%+v'", k, v))
		default:
			filterClauseWhere = append(filterClauseWhere, fmt.Sprintf("%s = %+v", k, v))
		}
	}
	if len(filterClauseWhere) > 0 {
		queryString += fmt.Sprintf(" WHERE %s", strings.Join(filterClauseWhere, " AND "))
	}

	queryString += " ;"
	_, err := d.sqlSvc.Exec(queryString)
	if err != nil {
		return err
	}
	return nil

}
