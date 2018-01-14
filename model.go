package main

type standard struct {
	URL       string  `json:"url"`
	PartNum   string  `json:"part"`
	Brand     string  `json:"brand"`
	Rank      int     `json:"rank"`
	Benchmark float32 `json:"benchmark"`
	Samples   int     `json:"samples"`
	Model     string  `json:"model"`
}

type cpu struct {
	Cores       string    `json:"cores"`
	Scores      [3]string `json:"scores"`
	SegmentPerf [3]string `json:"performance"`
	SubResults  [9]string `json:"subresults"`
	standard
}

type gpu struct {
	//lighting, reflection, parallax
	//mrender, gravity, splatting
	Name       string
	SubResults [6]string
	Averages   [2]string
}

type ssd struct {
	Name, Controller string
	SubResults       [9]string
	Averages         [3]string
	standard
}

var ssds []ssd

var gpus []gpu

var cpus = make(map[string]cpu)
