// config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser      string
	DBPassword  string
	DBName      string
	DBHost      string
	DBPort      int
	RedisHost   string
	RedisPort   int
	LogFilePath string
	LogFileName string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config := &Config{
		DBUser:      viper.GetString("DB_USER"),
		DBPassword:  viper.GetString("DB_PASSWORD"),
		DBName:      viper.GetString("DB_NAME"),
		DBHost:      viper.GetString("DB_HOST"),
		DBPort:      viper.GetInt("DB_PORT"),
		RedisHost:   viper.GetString("REDIS_HOST"),
		RedisPort:   viper.GetInt("REDIS_PORT"),
		LogFilePath: viper.GetString("LOG_FILE_PATH"),
		LogFileName: viper.GetString("LOG_FILE_NAME"),
	}

	return config
}
