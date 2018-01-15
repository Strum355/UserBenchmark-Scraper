package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func (c cpu) getCPU(ctx context.Context, cdp *chromedp.CDP, url string) (cpu, error) {
	var html string

	if err := cdp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(c.URL),
		chromedp.OuterHTML(`html`, &html, chromedp.ByQuery),
	}); err != nil {
		fmt.Printf("%s %s\n", c.URL, err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return c, err
	}

	//Cores
	c.Cores = doc.Find(`.cmp-cpt.tallp.cmp-cpt-l`).Text()

	//Scores - averages
	for i := 0; i < 3; i++ {
		c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3)).Text()
		if c.Scores[i] == "" {
			c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3)).Text()
			if c.Scores[i] == "" {
				c.Scores[i] = doc.Find(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3)).Text()
			}
		}
	}

	//SubResults - pre averages
	doc.Find(`.mcs-hl-col`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		c.SubResults[i] = s.Text()
		if i == 8 {
			return true
		}
		return false
	})

	//SegmentPerf - Gaming, Desktop, Workstation
	doc.Find(`.bsc-w.text-left.semi-strong`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		c.SegmentPerf[i] = s.Text()
		if i == 2 {
			return true
		}
		return false
	})

	fmt.Println(&c)

	return c, nil
}
