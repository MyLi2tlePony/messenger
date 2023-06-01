package server

import (
	"net"

	"github.com/spf13/viper"
)

type config struct {
	Host string
	Port string
}

func New(configPath string) (*config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		Host: viper.GetString("server.Host"),
		Port: viper.GetString("server.Port"),
	}, nil
}

func (s *config) GetPort() string {
	return s.Port
}

func (s *config) GetHost() string {
	return s.Host
}

func (s *config) GetHostPort() string {
	return net.JoinHostPort(s.Host, s.Port)
}
