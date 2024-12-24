//Viper的详细用法
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}

	ONNX struct {
		ModelPath string
		UseGPU    bool
	}

	Server struct {
		Port string
	}
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	return config
}
