package logger

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type config struct {
	Level string
}

func New(configPath string) (*config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		Level: viper.GetString("logger.level"),
	}, nil
}

func (l config) GetLevel() log.Lvl {
	switch l.Level {
	case "info":
		return log.INFO
	case "debug":
		return log.DEBUG
	case "error":
		return log.ERROR
	case "warn":
		return log.WARN
	}

	return log.OFF
}
