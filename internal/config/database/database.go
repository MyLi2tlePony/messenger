package database

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	Prefix       string
	DatabaseName string
	Host         string
	Port         string
	UserName     string
	Password     string
}

func New(configPath string) (*config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		Prefix:       viper.GetString("database.Prefix"),
		DatabaseName: viper.GetString("database.DatabaseName"),
		Host:         viper.GetString("database.Host"),
		Port:         viper.GetString("database.Port"),
		UserName:     viper.GetString("database.UserName"),
		Password:     viper.GetString("database.Password"),
	}, nil
}

func (d *config) GetPrefix() string {
	return d.Prefix
}

func (d *config) GetDatabaseName() string {
	return d.DatabaseName
}

func (d *config) GetHost() string {
	return d.Host
}

func (d *config) GetPort() string {
	return d.Port
}

func (d *config) GetUserName() string {
	return d.UserName
}

func (d *config) GetPassword() string {
	return d.Password
}

func (d *config) GetConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		d.Prefix, d.UserName, d.Password, d.Host, d.Port, d.DatabaseName)
}
