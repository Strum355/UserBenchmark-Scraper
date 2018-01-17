package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func (c CPU) getCPU(ctx context.Context, cdp *chromedp.CDP, url string) (CPU, error) {
	var html string

	getOuterHTML(ctx, c, cdp, &html)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return c, err
	}

	c.getCores(*doc)
	c.getAverages(*doc)
	c.getSubResults(*doc)
	c.getRelPerf(*doc)

	fmt.Println(&c)

	if old, ok := cpus.get(c.Model); ok {
		if !c.isValid(old) {
			return c, ErrNotValid
		}
	}

	return c, nil
}

func (c *CPU) getCores(doc goquery.Document) {
	c.Cores = doc.Find(`.cmp-cpt.tallp.cmp-cpt-l`).Text()
}

func (c *CPU) getAverages(doc goquery.Document) {
	for i := 0; i < 3; i++ {
		c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3)).Text()
		if c.Scores[i] == "" {
			c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3)).Text()
			if c.Scores[i] == "" {
				c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3)).Text()
			}
		}
	}
}

func (c *CPU) getSubResults(doc goquery.Document) {
	doc.Find(`.mcs-hl-col`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		c.SubResults[i] = s.Text()
		if i == 8 {
			return true
		}
		return false
	})
}

func (c *CPU) getRelPerf(doc goquery.Document) {
	doc.Find(`.bsc-w.text-left.semi-strong`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		c.SegmentPerf[i] = s.Text()
		if i == 2 {
			return true
		}
		return false
	})
}

func (c CPU) isValid(old CPU) bool {
	switch {
	case c.Cores == "" && old.Cores != "":
		return false
	case !equallyEmpty(c.Scores[:], old.Scores[:]):
		return false
	case !equallyEmpty(c.SegmentPerf[:], old.SegmentPerf[:]):
		return false
	case !equallyEmpty(c.SubResults[:], old.SubResults[:]):
		return false
	default:
		return true
	}
}

func equallyEmpty(new, old []string) bool {
	for i, v := range new {
		if v == "" && old[i] != "" {
			return false
		}
	}
	return true
}
