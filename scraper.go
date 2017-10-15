package main

import (
	"sync"
	"io/ioutil"
	"encoding/json"
	"time"
	"strings"
	"context"
	"fmt"
	"log"
	"os"
	"encoding/csv"

	cdp "github.com/knq/chromedp"
	cdptypes "github.com/knq/chromedp/cdp"
)

type result struct {
	Cores, Name string
	Scores, SegmentPerf [3]string
}

var results []result
var oldResults []result
var c *cdp.CDP

func main() {
	f, err := os.OpenFile("CPU_DATA.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file. Empty?", err)	
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&results)
	if err != nil {
		fmt.Println("Couldnt decode", err )
	}

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err = cdp.New(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	text(ctxt)

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("JSON: Your time was just wasted :( ", err)
		return
	}
	
	err = ioutil.WriteFile("CPU_DATA.json", out, 0644)
	if err != nil {
		fmt.Println("WRITE: Your time was just wasted :( ", err)		
	}
}

func text(ctx context.Context) {
	fmt.Println("STARTING")
	if err := c.Run(ctx, cdp.Tasks{
		cdp.Navigate(`http://www.userbenchmark.com/page/login`),
		cdp.SetValue(`input[name="username"]`, ""),
		cdp.Sleep(time.Second*2),
		cdp.SetValue(`input[name="password"]`, ""),
		cdp.Sleep(time.Second*2),
		cdp.Click(`button[name="submit"]`),
		cdp.Sleep(time.Second*3),
	}); err != nil {
		fmt.Println(err)
	}

	for name, val := range parseCSV("CPU_UserBenchmarks.csv") {
		if isIn(name) {
			fmt.Println("Skipping entry already saved")
			continue
		}
		var res result
		res.Name = name
		getPage(ctx, &res, val)
		out, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Println("JSON: ", err)
			return
		}
		
		err = ioutil.WriteFile("CPU_DATA.json", out, 0644)
		if err != nil {
			fmt.Println("WRITE: ", err)		
		}
	}
}

func isIn(name string) bool {
	for _, r := range results {
		if r.Name == name {
			return true
		}
	}

	return false
}

func getPage(ctxt context.Context, res *result, url string) {
	fmt.Println("Going to ", url)
	if err := c.Run(ctxt, cdp.Navigate(url)); err != nil {
		fmt.Println("Error navigating to ", url, err)
	}

	fmt.Println("Trying cores")
	if err := c.Run(ctxt, cdp.Text(`.cmp-cpt.tallp.cmp-cpt-l`, &res.Cores)); err != nil {
		fmt.Println("Error getting cores for ", url, err)
	}

	var wait sync.WaitGroup
	for i := 0; i < 3; i++ {
		wait.Add(2)
		
		go func(i int) {
			fmt.Println(i, "Trying green scores")		
			err := func() error {
				ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
				defer cancel()
				if err := c.Run(ctxt, cdp.Text(`.mcs-caption.pgbg`, &res.Scores[i])); err != nil {
					fmt.Println(i, "Error getting scores for ", url, err, "\ntrying yellow")
					return err
				}
				if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-caption.pgbg`, "class", "")); err != nil {
					fmt.Println(i, ".mcs-caption.pgbg", url, err)
				}
				return nil
			}()

			//if green fails, try yellow
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(`.mcs-caption.pybg`, &res.Scores[i])); err != nil {
						fmt.Println(i, "Error getting scores for ", url, err, "\ntrying red")
						return err
					}
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-caption.pybg`, "class", "")); err != nil {
						fmt.Println(i, ".mcs-caption.pybg", url, err)
					}
					return nil
				}()
			}else{
				wait.Done()
				fmt.Println(i, "found green")
				return
			}

			//if yellow fails, try red
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(`.mcs-caption.prbg`, &res.Scores[i])); err != nil {
						fmt.Println(i, "IMPORTANT!! .mcs-caption.prbg Error getting scores for ", url, err)
						return err
					}
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-caption.prbg`, "class", "")); err != nil {
						fmt.Println(i, "IMPORTANT!! .mcs-caption.prbg", url, err)
					}
					return nil
				}()
			}else{
				wait.Done()
				fmt.Println(i, "found yellow")
				return
			}

			if err == nil {
				fmt.Println("found red")
			}
			wait.Done()
		}(i)
		go func(i int) {
			fmt.Println("Trying performance ", i)
			if err := c.Run(ctxt, cdp.Text(`.bsc-w.text-left.semi-strong`, &res.SegmentPerf[i])); err != nil {
				fmt.Println(i, "Error getting performance for ", url, err)
			}
			if err := c.Run(ctxt, cdp.SetAttributeValue(`.bsc-w.text-left.semi-strong`, "class", "")); err != nil {
				fmt.Println(i, url, err)
			}
			wait.Done()
		}(i)
		
		wait.Wait()
	}

	if err := c.Run(ctxt, cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
		for i, val := range res.SegmentPerf {
			res.SegmentPerf[i] = strings.Trim(strings.Replace(strings.Replace(val, "\t", "", -1), "\n", " ", -1), " ")
		}
		return nil
	})); err != nil {
		fmt.Println(url, err)
	}
	
	if err := c.Run(ctxt, cdp.ActionFunc(func(context.Context, cdptypes.Handler) error {
		results = append(results, *res)
		return nil
	})); err != nil {
		fmt.Println(url, err)
	}
}

func parseCSV(filename string) (out map[string]string) {
	out = make(map[string]string)
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	reader := csv.NewReader(file)
	columns, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for i := 1; i < len(columns); i++ {
		if _, ok := out[columns[i][3]]; ok {
			continue
		}
		out[columns[i][3]] = columns[i][7]
	}

	return
}