package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func SetConfigFile(name, path, extension string) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.SetConfigType(extension)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", err)
	}
}

func GetString(key string) string {
	return viper.GetString(fmt.Sprintf("%v", key))
}
func GetInt(key string) int {
	return viper.GetInt(fmt.Sprint(key))
}
