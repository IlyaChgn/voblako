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
	DBName   string `env:"POSTGRES_DB"`
	Port     string `env:"POSTGRES_PORT"`
}

type PostgresAuthConfig struct {
	PostgresConfig
	Host string `env:"POSTGRES_AUTH_HOST"`
}

type PostgresFileConfig struct {
	PostgresConfig
	Host string `env:"POSTGRES_FILE_HOST"`
}

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	DB       int    `env:"REDIS_DB"`
}

type MinioConfig struct {
	Host      string `env:"MINIO_HOST"`
	Port      string `env:"MINIO_API_PORT"`
	AccessKey string `env:"MINIO_ROOT_USER"`
	SecretKey string `env:"MINIO_ROOT_PASSWORD"`
	Bucket    string `env:"MINIO_BUCKET"`
}

type AuthServiceConfig struct {
	Postgres PostgresAuthConfig
	Redis    RedisConfig

	InternalHost string `yaml:"host"`
	ExternalHost string `env:"AUTH_HOST"`
	Port         string `env:"AUTH_PORT"`
}

type FileServiceConfig struct {
	Postgres PostgresFileConfig
	Minio    MinioConfig

	InternalHost string `yaml:"host"`
	ExternalHost string `env:"FILE_HOST"`
	Port         string `env:"FILE_PORT"`
}

type CtxKeys struct {
	User string `yaml:"user"`
}

type Config struct {
	Server ServerConfig      `yaml:"server"`
	Auth   AuthServiceConfig `yaml:"auth_service"`
	File   FileServiceConfig `yaml:"file_service"`

	Keys CtxKeys `yaml:"ctx_keys"`
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
