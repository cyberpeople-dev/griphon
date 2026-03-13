package internal

import (
	"errors"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Validation
	if len(cfg.Endpoints) == 0 {
		return nil, errors.New("config: endpoints must not be empty")
	}
	if len(cfg.Payloads) == 0 {
		return nil, errors.New("config: payloads must not be empty")
	}
	if len(cfg.UserAgents) == 0 {
		return nil, errors.New("config: user_agents must not be empty")
	}
	return &cfg, nil
} 