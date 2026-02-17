package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	NATS     NATSConfig
	Influx   InfluxConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type NATSConfig struct {
	URL string
}

type InfluxConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type ServerConfig struct {
	Port int
	Mode string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("server.port", 8008)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.db", 0)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
