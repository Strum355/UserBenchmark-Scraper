package component

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
)

const (
	// Cores
	cores = `.cmp-cpt.tallp.cmp-cpt-l`
	// Subresults
	subResults = `.mcs-hl-col`
	// Averages
	firstAverage  = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`
	secondAverage = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`
	thirdAverage  = `.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`
	// Relative Performance
	firstRelative  = `.bsc-w.text-left.semi-strong div:first-child`
	secondRelative = `.bsc-w.text-left.semi-strong div:nth-child(3)`
)

type CPU struct {
	Cores      string    `json:"cores,omitempty"` //Cores
	Averages   [3]string `json:"averages"`        //Averages
	RelPerf    [3]string `json:"performance"`     //Relative Performance
	SubResults [9]string `json:"subresults"`      //Sub Results
	Standard
}

type CPUs struct {
	sync.RWMutex
	norm map[string]CPU
	rank map[int]CPU
}

func (c CPU) Get(ctx context.Context, cdp *cdp.Client) (Component, error) {
	html, err := GetOuterHTML(ctx, c, cdp)
	if err != nil {
		fmt.Println(err)
		return c, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html.OuterHTML))
	if err != nil {
		fmt.Println(err)
		return c, err
	}

	c.GetCores(doc)
	c.GetAverages(doc)
	c.GetSubResults(doc)
	c.GetRelativePerf(doc)

	/* 	if old, ok := cpus.Get(c.Model); ok {
		return c, c.IsValid(old)
	} */

	return c, nil
}

func (c CPU) GetURL() string {
	return c.URL
}

func (c *CPU) GetCores(doc *goquery.Document) {
	c.Cores = doc.Find(cores).Text()
}

func (c *CPU) GetAverages(doc *goquery.Document) {
	for i := 0; i < 3; i++ {
		c.Averages[i] = doc.Find(fmt.Sprintf(firstAverage, i+3)).Text()
		if c.Averages[i] == "" {
			c.Averages[i] = doc.Find(fmt.Sprintf(secondAverage, i+3)).Text()
			if c.Averages[i] == "" {
				c.Averages[i] = doc.Find(fmt.Sprintf(thirdAverage, i+3)).Text()
			}
		}
	}
}

func (c *CPU) GetSubResults(doc *goquery.Document) {
	doc.Find(subResults).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 8 {
			return false
		}
		c.SubResults[i] = s.Text()
		return true
	})
}

func (c *CPU) GetRelativePerf(doc *goquery.Document) {
	doc.Find(firstRelative).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 3 {
			return false
		}
		c.RelPerf[i] = strings.TrimSpace(s.Text())
		return true
	})

	doc.Find(secondRelative).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 3 {
			return false
		}
		c.RelPerf[i] += " " + s.Text()
		return true
	})
}

func (c CPU) IsValid(old Component) error {
	return nil
}