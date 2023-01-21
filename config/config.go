package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	rawConfig map[string]interface{}
}

func Create(configFile string) (*Config, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(b, &cfg.rawConfig)
	if err != nil {
		return nil, err
	}

	fmt.Println(cfg)

	return cfg, nil
}

func MustCreate(configFile string) Config {
	cfg, err := Create(configFile)
	if err != nil {
		panic(err)
	}

	return *cfg
}

func (c *Config) MustLoad(key string, to interface{}) {
	if err := c.Load(key, to); err != nil {
		panic(err)
	}
}

func (c *Config) Load(key string, to interface{}) error {
	b, err := yaml.Marshal(c.rawConfig[key])
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, to)
}
