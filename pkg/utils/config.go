package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	viper *viper.Viper
}

var configInstance *Config // package private singleton instance of the configuration
var singleton sync.Once    // package private singleton helper utility

func GetConfig() *viper.Viper {
	// create an instance if not available
	singleton.Do(func() {
		configInstance = &Config{viper.New()}
	})

	return configInstance.viper
}

func Start() {
	if configInstance == nil {
		GetConfig()
	}

	// Find and read the config file
	err := configInstance.viper.ReadInConfig()
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}
