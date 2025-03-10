package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string        `yaml:"env" env-defolt:"local"`
	ConnectionString string        `yaml:"connection_string" env-required:"./data"`
	TokenTTL         time.Duration `yaml:"token_ttl" env-required:"true"`
	Server           ServerConfig  `yaml:"server"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("Config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(cofigPath string) *Config {

	if _, err := os.Stat(cofigPath); os.IsNotExist(err) {
		panic("Config file doesn't exist: " + cofigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cofigPath, &cfg); err != nil {
		panic("failed to read config" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
