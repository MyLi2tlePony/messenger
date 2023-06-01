package database

import (
	"os"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Database *config
}

func Test(t *testing.T) {
	t.Run("read config case", func(t *testing.T) {
		fileName := "testConfig.*.toml"
		file, err := os.CreateTemp("", fileName)
		require.Nil(t, err)

		defer func() {
			require.Nil(t, file.Close())
			require.Nil(t, os.Remove(file.Name()))
		}()

		expectedConfig := &Config{
			Database: &config{
				Prefix:       "postgresql",
				DatabaseName: "postgres",
				Host:         "localhost",
				Port:         "5432",
				UserName:     "postgres",
				Password:     "1234512345",
			},
		}

		marshal, err := toml.Marshal(expectedConfig)
		require.Nil(t, err)

		_, err = file.Write(marshal)
		require.Nil(t, err)

		config, err := New(file.Name())
		require.Nil(t, err)

		require.True(t, reflect.DeepEqual(config, expectedConfig.Database))
	})
}
