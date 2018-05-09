package main

import (
	"sync"
)

type (
	Component interface {
		//Get(context.Context, *chromedp.CDP, string) (Component, error)
		GetURL() string
	}

	Standard struct {
		URL       string  `json:"url,omitempty"`
		PartNum   string  `json:"part,omitempty"`
		Brand     string  `json:"brand,omityempty"`
		Rank      int     `json:"rank,omityempty"`
		Benchmark float32 `json:"benchmark,omitempty"`
		Samples   int     `json:"samples,omitempty"`
		Model     string  `json:"model,omitempty"`
	}

	CPU struct {
		Cores       string    `json:"cores,omitempty"` //Cores
		Averages    [3]string `json:"scores"`          //Averages
		SegmentPerf [3]string `json:"performance"`     //Relative Performance
		SubResults  [9]string `json:"subresults"`      //Sub Results
		Standard
	}

	CPUs struct {
		sync.RWMutex
		norm map[string]CPU
		rank map[int]CPU
	}

	GPU struct {
		SubResults [6]string `json:"subresults"`
		Averages   [2]string `json:"scores"`
		Standard
	}

	GPUs struct {
		sync.RWMutex
		norm map[string]GPU
		rank map[int]GPU
	}

	SSD struct {
		Controller string    `json:"controller"`
		SubResults [9]string `json:"subresults"`
		Averages   [3]string `json:"scores"`
		Standard
	}

	SSDs struct {
		sync.RWMutex
		norm map[string]SSD
		rank map[int]SSD
	}
)

func (c *CPUs) Get(in string) (CPU, bool) {
	c.RLock()
	defer c.RUnlock()
	ret, ok := c.norm[in]
	return ret, ok
}

func (c *CPUs) Set(in string, v CPU) {
	c.Lock()
	defer c.Unlock()
	c.norm[in] = v
}

func (g *GPUs) Get(in string) (GPU, bool) {
	g.RLock()
	defer g.RUnlock()
	ret, ok := g.norm[in]
	return ret, ok
}

func (g *GPUs) Set(in string, v GPU) {
	g.Lock()
	defer g.Unlock()
	g.norm[in] = v
}

func (s *SSDs) Get(in string) (SSD, bool) {
	s.RLock()
	defer s.RUnlock()
	ret, ok := s.norm[in]
	return ret, ok
}

func (s *SSDs) Set(in string, v SSD) {
	s.Lock()
	defer s.Unlock()
	s.norm[in] = v
}
