package main

import (
	"sync"
)

type (
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

	gpu struct {
		//lighting, reflection, parallax
		//mrender, gravity, splatting
		Name       string
		SubResults [6]string
		Averages   [2]string
	}

	ssd struct {
		Name, Controller string
		SubResults       [9]string
		Averages         [3]string
		Standard
	}

	CPUs struct {
		*sync.RWMutex
		c map[string]CPU
	}
)

func (c *CPUs) get(s string) (CPU, bool) {
	c.RLock()
	ret, ok := c.c[s]
	c.RUnlock()
	return ret, ok
}

func (c *CPUs) set(s string, v CPU) {
	c.Lock()
	c.c[s] = v
	c.Unlock()
}
