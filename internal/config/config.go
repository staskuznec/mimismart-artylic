package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel string `yaml:"logLevel"`
}

func LoadConfig(path string) Config {
	cfg := Config{LogLevel: "error"}
	file, err := os.Open(path)
	if err != nil {
		return cfg
	}
	defer file.Close()
	_ = yaml.NewDecoder(file).Decode(&cfg)
	cfg.LogLevel = strings.ToLower(cfg.LogLevel)
	return cfg
}
