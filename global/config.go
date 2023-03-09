package global

import (
	cfg "github.com/mailio/go-web3-kit/config"
)

// Conf global config
var Conf Config

type Config struct {
	cfg.YamlConfig `yaml:",inline"`
	CouchDB        CouchDBConfig `yaml:"couchdb"`
}

type CouchDBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
