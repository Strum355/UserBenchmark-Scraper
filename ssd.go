package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func (s SSD) Get(ctx context.Context, cdp *chromedp.CDP, url string) (Component, error) {
	var html string

	GetOuterHTML(ctx, s, cdp, &html)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return s, err
	}

	s.GetSubResults(*doc)
	s.GetAverages(*doc)

	fmt.Println(&s)

	if old, ok := ssds.Get(s.Model); ok {
		if !s.IsValid(old) {
			return s, ErrNotValid
		}
	}

	return s, nil
}

func (s SSD) GetURL() string {
	return s.URL
}

func (s *SSD) GetSubResults(doc goquery.Document) {
	doc.Find(`mcs-hl-col`).EachWithBreak(func(i int, sel *goquery.Selection) bool {
		if i == 6 {
			return true
		}
		s.SubResults[i] = sel.Text()
		return false
	})
}

func (s *SSD) GetAverages(doc goquery.Document) {
	for i := 0; i < 2; i++ {
		s.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3)).Text()
		if s.Averages[i] == "" {
			s.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3)).Text()
			if s.Averages[i] == "" {
				s.Averages[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3)).Text()
			}
		}
	}
}

func (s SSD) IsValid(old SSD) bool {
	switch {
	case !equallyEmpty(s.Averages[:], old.Averages[:]):
		return false
	case !equallyEmpty(s.SubResults[:], old.SubResults[:]):
		return false
	case s.Controller != old.Controller:
		return false
	default:
		return true
	}
}
