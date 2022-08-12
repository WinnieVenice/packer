package conf

import (
	"path"
	"runtime"

	"github.com/spf13/viper"
)

/**
import (
	_ "github.com/WinnieVenice/packer/conf"
)
*/
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	configPath := path.Dir(path.Dir(fileName))
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
