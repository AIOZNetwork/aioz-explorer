package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var viperGoB *viper.Viper
func LoadConfig(configPath string) {
	viperGoB = viper.New()
	viperGoB.SetConfigFile(configPath)

	if err := viperGoB.ReadInConfig(); err != nil {
		fmt.Println("Error loading config.toml", err.Error())
		panic("Error loading config.toml")
	}

	viperGoB.WatchConfig()
	viperGoB.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

func GetConfig() *viper.Viper {
	return viperGoB
}