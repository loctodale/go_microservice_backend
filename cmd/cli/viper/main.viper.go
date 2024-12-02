package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"PORT"`
	} `mapstructure:"SERVER"`
	Database []struct {
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Host     string `mapstructure:"HOST"`
	} `mapstructure:"DATABASE"`
}

func main() {
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
	fmt.Println("Server port::", viper.GetInt("server.port"))
	fmt.Println("Jwt Key::", viper.GetString("security.jwt.key"))

	//config structure
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode config %v", err)
	}

	fmt.Println("Config Server port::", config.Server.Port)

	for _, db := range config.Database {
		fmt.Printf("dabase user: %s , password: %s, hostname: %s \n", db.User, db.Password, db.Host)
	}
}
