package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var configurations *Configuration

type DatabaseVariable struct {
	User     string
	DB_Name  string
	Host     string
	Password string
	Port     int
}

type Configuration struct {
	Database DatabaseVariable
	Log      LoggingConfiguration
}

type LoggingConfiguration struct {
	Level string
}

func InitConfig() {
	configurations = new(Configuration)
}

func GetConfig() *Configuration {
	return configurations
}
func InitConfiguration() {
	viper.SetConfigName("config.yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file " + err.Error())
	}
	err := viper.Unmarshal(&configurations)
	fmt.Println("Configuration", configurations)
	if err != nil {
		panic("Not Able To Read The File")
	}
}
