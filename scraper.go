package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Strum355/UserBenchmark-Scraper/config"

	"github.com/Strum355/UserBenchmark-Scraper/chrome"
)

var (
	conf = new(config.Config)
)

func init() {
	conf.LoadConfig()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*4))
	defer cancel()

	c, conn, err := chrome.Start(ctx, "http://127.0.0.1:9222")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	if err := chrome.Login(ctx, c, conf); err != nil {
		fmt.Println(err)
		return
	}
}
