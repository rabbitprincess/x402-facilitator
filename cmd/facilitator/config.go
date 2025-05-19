package main

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rabbitprincess/x402-facilitator/types"
)

type Config struct {
	Scheme     types.Scheme `mapstructure:"scheme"`
	Network    string       `mapstructure:"network"`
	Port       int          `mapstructure:"port"`
	Url        string       `mapstructure:"url"`
	PrivateKey string       `mapstructure:"private_key"`
}

func LoadConfig(path string) (*Config, error) {
	var k = koanf.New(".")

	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, err
	}
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
