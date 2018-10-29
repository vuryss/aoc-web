package helper

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	initialized bool
	file string
	configuration map[string]interface{}
}

func (config *Config) Get(key string) string {
	return ""
}

func (config *Config) getConfig() {
	if !config.initialized {
		config.file = ResolveProjectFile("config/app-config.json")

		content, err := ioutil.ReadFile(config.file)

		if err != nil {
			panic("Cannot read project configuration file: " + config.file)
		}

		err = json.Unmarshal(content, &config.configuration)

		if err != nil {
			panic("Cannot decode project configuration file: " + config.file)
		}
	}
}
