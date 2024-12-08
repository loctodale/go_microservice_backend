package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go_microservice_backend_api/global"
)

func LoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// read server configuration
	//fmt.Println("Server port::", viper.GetInt("server.port"))
	//fmt.Println("Jwt Key::", viper.GetString("security.jwt.certs"))

	//config structure
	if err = viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("Unable to decode config %v", err)
	}

	//fmt.Println("Config Server port::", global.Config.Server.Port)
	//
	//for _, db := range config.Database {
	//	fmt.Printf("dabase user: %s , password: %s, hostname: %s \n", db.User, db.Password, db.Host)
	//}

}
