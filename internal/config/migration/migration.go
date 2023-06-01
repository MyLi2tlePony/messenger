package migration

import "github.com/spf13/viper"

type config struct {
	UpPath   string
	DownPath string
}

func New(configPath string) (*config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		UpPath:   viper.GetString("migration.upPath"),
		DownPath: viper.GetString("migration.downPath"),
	}, nil
}

func (l config) GetUpPath() string {
	return l.UpPath
}

func (l config) GetDownPath() string {
	return l.DownPath
}
