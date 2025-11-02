package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host    string   `yaml:"host"`
	Port    string   `env:"APP_PORT"`
	Timeout int      `yaml:"timeout"`
	Origins []string `yaml:"origins"`
	Headers []string `yaml:"headers"`
	Methods []string `yaml:"methods"`
}

type PostgresConfig struct {
	Username string `env:"POSTGRES_USERNAME"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_DB"`
}

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	DB       int    `env:"REDIS_DB"`
}

type ServiceConfig struct {
	InternalHost string `yaml:"host"`
	ExternalHost string `env:"AUTH_HOST"`
	Port         string `env:"AUTH_PORT"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`

	Postgres PostgresConfig
	Redis    RedisConfig
	Auth     ServiceConfig `yaml:"auth_service"`
}

func ReadConfig(cfgPath string) *Config {
	cfg := &Config{}

	file, err := os.Open(cfgPath)
	if err != nil {
		log.Println("Something went wrong while opening config file ", err)

		return nil
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		log.Println("Something went wrong while reading config from yaml file ", err)

		return nil
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Println("Something went wrong while reading config from env file ", err)

		return nil
	}

	log.Println("Successfully opened config")

	return cfg
}
