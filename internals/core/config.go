package core

import (
	"../helper"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func NewConfig(configFile string) *Config {
	file := helper.ResolveProjectFile(configFile)
	content, err := ioutil.ReadFile(file)

	if err != nil {
		panic("CORE CONFIG: Cannot read project configuration file: " + file)
	}

	configInstance := &Config{}

	err = json.Unmarshal(content, &configInstance.configuration)

	if err != nil {
		panic("CORE CONFIG: Cannot decode project configuration file: " + file)
	}

	return configInstance
}

type Config struct {
	configuration map[string]interface{}
}

func (config *Config) Get(key string) (interface{}, bool) {
	path := strings.Split(key, ".")
	data := config.configuration
	hasChildren := true

	var elem interface{}
	var exists bool

	for i := range path {
		if !hasChildren {
			return nil, false
		}

		elem, exists = data[path[i]]

		if !exists {
			return nil, false
		}

		data, hasChildren = elem.(map[string]interface{})
	}

	return elem, true
}

func (config *Config) GetString(key string) (string, bool) {
	value, exists := config.Get(key)
	stringValue, ok := value.(string)

	if !ok {
		return "", false
	}

	return stringValue, exists
}
