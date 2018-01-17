package main

import (
	"sync"
	"context"
	"github.com/chromedp/chromedp"
)

type (
	Component interface {
		Get(context.Context, *chromedp.CDP, string) (Component, error)
		GetURL() string
	}

	Standard struct {
		URL       string  `json:"url"`
		PartNum   string  `json:"part"`
		Brand     string  `json:"brand"`
		Rank      int     `json:"rank"`
		Benchmark float32 `json:"benchmark"`
		Samples   int     `json:"samples"`
		Model     string  `json:"model"`
	}

	CPU struct {
		Cores       string    `json:"cores"`       //Cores
		Scores      [3]string `json:"scores"`      //Averages
		SegmentPerf [3]string `json:"performance"` //Relative Performance
		SubResults  [9]string `json:"subresults"`  //Sub Results
		Standard
	}

	CPUs struct {
		sync.RWMutex
		norm map[string]*CPU
		rank map[int]*CPU
	}

	GPU struct {
		//lighting, reflection, parallax
		//mrender, gravity, splatting
		SubResults [6]string
		Averages   [2]string
		Standard
	}

	GPUs struct {
		sync.RWMutex
		norm map[string]*GPU
		rank map[int]*GPU
	}

	ssd struct {
		Name, Controller string
		SubResults       [9]string
		Averages         [3]string
		Standard
	}


)

/*
	if ret != nil {} may be premature optimization or not even an optimization
	at all, will look at later in life
*/

func (c *CPUs) Get(s string) (CPU, bool) {
	c.RLock()
	ret, ok := c.norm[s]
	c.RUnlock()
	if ret == nil {
		return CPU{}, ok
	}
	return *ret, ok
}

func (c *CPUs) Set(s string, v CPU) {
	c.Lock()
	c.norm[s] = &v
	c.Unlock()
}

func (g *GPUs) Get(s string) (GPU, bool) {
	g.RLock()
	ret, ok := g.norm[s]
	g.RUnlock()
	if ret == nil {
		return GPU{}, ok
	}
	return *ret, ok
}

func (g *GPUs) Set(s string, v GPU) {
	g.Lock()
	g.norm[s] = &v
	g.Unlock()
}