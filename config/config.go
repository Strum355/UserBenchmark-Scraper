package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (c *Config) LoadConfig() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if err := json.Unmarshal(b, c); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
