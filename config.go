package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (c *config) loadConfig() {
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
