package component

import (
	"context"
	"sync"

	"github.com/Strum355/UserBenchmark-Scraper/chrome"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
)

type Component interface {
	Get(context.Context, *cdp.Client) (Component, error)
	GetURL() string
	IsValid(Component) error
}

type Standard struct {
	URL       string  `json:"url,omitempty"`
	PartNum   string  `json:"part,omitempty"`
	Brand     string  `json:"brand,omityempty"`
	Rank      int     `json:"rank,omityempty"`
	Benchmark float32 `json:"benchmark,omitempty"`
	Samples   int     `json:"samples,omitempty"`
	Model     string  `json:"model,omitempty"`
}

type GPU struct {
	SubResults [6]string `json:"subresults"`
	Averages   [2]string `json:"averages"`
	Standard
}

type GPUs struct {
	sync.RWMutex
	norm map[string]GPU
	rank map[int]GPU
}

type SSD struct {
	Controller string    `json:"controller"`
	SubResults [9]string `json:"subresults"`
	Averages   [3]string `json:"averages"`
	Standard
}

type SSDs struct {
	sync.RWMutex
	norm map[string]SSD
	rank map[int]SSD
}

func GetOuterHTML(ctx context.Context, c Component, cdp *cdp.Client) (result *dom.GetOuterHTMLReply, err error) {
	if err = chrome.Navigate(ctx, cdp, c.GetURL()); err != nil {
		return
	}

	var doc *dom.GetDocumentReply
	doc, err = cdp.DOM.GetDocument(ctx, nil)
	if err != nil {
		return
	}

	result, err = cdp.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})
	if err != nil {
		return
	}

	return
}
