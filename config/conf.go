package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed conf.yaml
var confbytes []byte

type Conf struct {
	Channels  map[int]string `yaml:"channels"`
	MQTT      string         `yaml:"mqtt"`
	RadioPort string         `yaml:"radio_port"`
	HookDB    string         `yaml:"hook_db"`
	CacheDB   string         `yaml:"cache_db"`
}

func Config() Conf {
	c := Conf{}
	if err := yaml.Unmarshal(confbytes, &c); err != nil {
		log.Panic(err)
	}
	return c
}
