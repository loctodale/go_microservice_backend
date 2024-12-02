package initialize

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"time"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}
func InitMysql() {
	m := global.Config.Mysql
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	//var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.DbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		m.Username, m.Password, m.Host, m.Port, m.DbName)
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	checkErrorPanic(err, "InitMysql initialization error")
	global.Logger.Info("InitMysql success")
	global.Mdb = db

	SetPool()
	genTableDAO()
	migateTables()
}

func SetPool() {
	m := global.Config.Mysql
	sqlDb := global.Mdb

	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns))
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifeTime))
}

func migateTables() {
	//err := global.Mdb.AutoMigrate(
	////&po.User{},
	////&po.Role{},
	////&model.GoCrmUserV2{},
	//)
	//if err != nil {
	//	fmt.Println("Mysql error: %s ::", err)
	//}
}
func genTableDAO() {
	//g := gen.NewGenerator(gen.Config{
	//	OutPath: "./internal/model",
	//	Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	//})
	//
	//g.UseDB(global.Mdb)
	//g.GenerateModel("go_crm_user")
	//
	//// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.ApplyBasic(model.User{})
	//
	//// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// Generate the code
	//g.Execute()
}
