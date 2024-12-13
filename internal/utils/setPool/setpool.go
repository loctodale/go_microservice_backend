package setPool

import (
	"database/sql"
	"go_microservice_backend_api/global"
	"time"
)

func SetPool(sqlDb *sql.DB, connectMaxTime int, maxOpenConnection int, lifeTime int) {
	sqlDb.SetConnMaxIdleTime(time.Duration(global.Config.Mysql.UserTable.MaxIdleConns))
	sqlDb.SetMaxOpenConns(global.Config.Mysql.UserTable.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Duration(global.Config.Mysql.UserTable.MaxIdleConns))
}
