package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Common     Common      `json:"common"`
	Sensors    []sensor    `json:"sensors"`
	Collectors []collector `json:"collectors"`
}

type Common struct {
	HostName   string `json:"hostName"`
	Debug      bool   `json:"debug"`
	PrivateKey string `json:"privateKey"`
}

type sensor struct {
	Name    string   `json:"name"`
	Enable  bool     `json:"enable"`
	Metrics []string `json:"metrics"`
	Devices []string `json:"devices"`
}

type collector struct {
	Name       string                 `json:"name"`
	Enable     bool                   `json:"enable"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (config *Config) Read(filename *string) (err error) {
	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			err = fmt.Errorf("[FATAL] %s", recoveryMessage)
		}
	}()

	jsonData, err := os.ReadFile(*filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		panic(fmt.Errorf("%s: %v", *filename, err))
	}

	return nil
}
