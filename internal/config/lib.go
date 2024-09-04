package config

import "github.com/kelseyhightower/envconfig"

func LoadConfig() (config Config, err error) {
	err = envconfig.Process("racing_draft", &config)
	return config, err
}

type Config struct {
	APIBaseURL string `envconfig:"api_base_url"`

	DB     DBConf
	Server ServerConf
}

type DBConf struct {
	Host     string `default:"postgres"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Port     uint16 `default:"5432"`
	DBName   string `default:"postgres"`
	Region   string `default:""`
}

type ServerConf struct {
	Address string `default:"0.0.0.0"`
	Port    uint16 `default:"8080"`
}
