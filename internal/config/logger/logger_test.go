package logger

import (
	"os"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Logger *config
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
			Logger: &config{
				Level: "info",
			},
		}

		marshal, err := toml.Marshal(expectedConfig)
		require.Nil(t, err)

		_, err = file.Write(marshal)
		require.Nil(t, err)

		config, err := New(file.Name())
		require.Nil(t, err)

		require.True(t, reflect.DeepEqual(config, expectedConfig.Logger))
	})
}
