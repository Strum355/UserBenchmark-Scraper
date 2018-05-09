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
	cpus = new(CPUs)
	gpus = new(GPUs)
	ssds = new(SSDs)
)

func init() {
	conf.LoadConfig()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*50))
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

	time.Sleep(time.Second * 2)

	/*doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	} 

	result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})
	if err != nil {

		return
	}*/
}
