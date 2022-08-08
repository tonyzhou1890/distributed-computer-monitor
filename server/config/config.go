package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GConfig *Config

func init() {
	GConfig = &Config{}
}

func InitConfig(env string) (err error) {
	var file = "./config/" + env + ".yml"

	yamlConf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlConf, GConfig); err != nil {
		return err
	}

	return nil
}
