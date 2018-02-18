package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func (g GPU) Get(ctx context.Context, cdp *chromedp.CDP, url string) (Component, error) {
	var html string

	GetOuterHTML(ctx, g, cdp, &html)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return g, err
	}

	g.GetSubResults(*doc)
	g.GetAverages(*doc)

	fmt.Println(&g)

	if old, ok := gpus.Get(g.Model); ok {
		if !g.IsValid(old) {
			return g, ErrNotValid
		}
	}

	return g, nil
}

func (g GPU) GetURL() string {
	return g.URL
}

func (g *GPU) GetSubResults(doc goquery.Document) {
	doc.Find(`mcs-hl-col`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 6 {
			return true
		}
		g.SubResults[i] = s.Text()
		return false
	})
}

func (g *GPU) GetAverages(doc goquery.Document) {
	for i := 0; i < 2; i++ {
		g.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3)).Text()
		if g.Averages[i] == "" {
			g.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3)).Text()
			if g.Averages[i] == "" {
				g.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3)).Text()
			}
		}
	}
}

func (g GPU) IsValid(old GPU) bool {
	switch {
	case !equallyEmpty(g.Averages[:], old.Averages[:]):
		return false
	case !equallyEmpty(g.SubResults[:], old.SubResults[:]):
		return false
	default:
		return true
	}
}
