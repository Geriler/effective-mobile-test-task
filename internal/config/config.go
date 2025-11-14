package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger      Logger        `yaml:"logger" env-prefix:"LOGGER_"`
	TimeoutStop time.Duration `env:"TIMEOUT_STOP" yaml:"timeout_stop" env-default:"10s"`
	HTTP        Address       `yaml:"http" env-prefix:"HTTP_"`
	GRPC        Address       `yaml:"grpc" env-prefix:"GRPC_"`
	Database    Database      `yaml:"database" env-prefix:"DATABASE_"`
}

type Address struct {
	Host string `env:"HOST" yaml:"host" env-required:"true"`
	Port int    `env:"PORT" yaml:"port" env-required:"true"`
}

type Database struct {
	Host     string `env:"HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"PORT" yaml:"port" env-required:"true"`
	User     string `env:"USER" yaml:"user" env-required:"true"`
	Password string `env:"PASSWORD" yaml:"password" env-required:"true"`
	Name     string `env:"NAME" yaml:"name" env-required:"true"`
}

type Logger struct {
	Type  string `env:"TYPE" yaml:"type" env-default:"json"`
	Level string `env:"LEVEL" yaml:"level" env-default:"info"`
}

func MustLoad() Config {
	var path string
	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	var config Config

	if path == "" {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			panic(err)
		}

		return config
	}

	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic(err)
	}

	return config
}
