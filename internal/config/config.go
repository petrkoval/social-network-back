package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Server   *ServerConfig `yaml:"server"`
	Database *DBConfig     `yaml:"database"`
	Tokens   *TokensConfig `yaml:"tokens"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type TokensConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

func MustLoad() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
