package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type Config struct {
	ListenAddr   string `yaml:"listen_addr"`
	UpstreamAddr string `yaml:"upstream_addr"`
}
