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

	// Explicitly bind environment variables for nested config keys.
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")
	_ = viper.BindEnv("database.host", "DATABASE_HOST")
	_ = viper.BindEnv("database.port", "DATABASE_PORT")
	_ = viper.BindEnv("database.user", "DATABASE_USER")
	_ = viper.BindEnv("database.password", "DATABASE_PASSWORD")
	_ = viper.BindEnv("database.dbname", "DATABASE_DBNAME")
	_ = viper.BindEnv("database.sslmode", "DATABASE_SSLMODE")
	_ = viper.BindEnv("redis.addr", "REDIS_ADDR")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "REDIS_DB")
	_ = viper.BindEnv("nats.url", "NATS_URL")
	_ = viper.BindEnv("influx.url", "INFLUX_URL")
	_ = viper.BindEnv("influx.token", "INFLUX_TOKEN")
	_ = viper.BindEnv("influx.org", "INFLUX_ORG")
	_ = viper.BindEnv("influx.bucket", "INFLUX_BUCKET")

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
