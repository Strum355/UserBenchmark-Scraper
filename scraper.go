package main

import (
	"errors"
	"context"
	_"encoding/csv"
	"fmt"
	"log"
	_"os"
	_"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	_ "github.com/fatih/color"
)

var (
	conf = new(config)
	cpus = new(CPUs)
	gpus = new(GPUs)
	// ErrNotValid is returned if fields of the item are missing fields compared to the previously 
	// stored entry.
	ErrNotValid = errors.New("fields were missed")
	// ErrNewEntry is returned if the item does not reside in the map.
	// This can be ignored if it is being scraped for the first time eg if its newly added to the CSV.
	ErrNewEntry = errors.New("entry not in map")
)

func main() {
	conf.loadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := chromedp.New(ctx, chromedp.WithRunnerOptions(
		runner.HeadlessPathPort("/usr/bin/google-chrome-stable", 9222),
		runner.Flag("headless", false)),
		chromedp.WithErrorf(func(s string, v ...interface{}) {
			if strings.Contains(s, "could not perform") || strings.Contains(s, "could not get") {
				return
			}
			log.Printf("error: "+s, v...)
		}))
	if err != nil {
		fmt.Println(err)
		return
	}

	start(ctx, c)

	if err := c.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func start(ctx context.Context, c *chromedp.CDP) {
	fmt.Println("STARTING")
	if err := login(ctx, c); err != nil {
		fmt.Println("Couldnt login", err)
		return
	}

	g := CPU{
		Standard: Standard{
			URL: "http://cpu.userbenchmark.com/Intel-Core-i3-8350K/Rating/3935",
		},
	}
	g.Get(ctx, c, g.URL)
	channel := make(chan struct{})
	<-channel
}

func login(ctx context.Context, cdp *chromedp.CDP) error {
	return cdp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(`http://www.userbenchmark.com/page/login`),
		chromedp.SetValue(`input[name="username"]`, conf.User),
		chromedp.SetValue(`input[name="password"]`, conf.Pass),
		chromedp.Sleep(time.Second),
		chromedp.Click(`button[name="submit"]`),
		chromedp.Sleep(time.Second*5),
	})
}

func GetOuterHTML(ctx context.Context, c Component, cdp *chromedp.CDP, s *string) error {
	return cdp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(c.GetURL()),
		chromedp.Sleep(time.Second*5),
		chromedp.WaitReady(`body`, chromedp.BySearch),
		chromedp.OuterHTML(`body`, s, chromedp.BySearch),
	})
}

/*
func doCPU(ctx context.Context, c *chromedp.CDP) {
	bytes, err := ioutil.ReadFile("./CPU_DATA_MAP.json")
	if err != nil {
		log.Fatalln(err)
	}

	in := make(map[string]cpu)
	err = json.Unmarshal(bytes, &in)
	if err != nil {
		log.Fatalln(err)
	}

Outer:
	for _, val := range parseCSV("CPU_UserBenchmarks.csv") {
		res := cpu{
			standard: standard{
				URL:       val.URL,
				PartNum:   val.PartNum,
				Brand:     val.Brand,
				Rank:      val.Rank,
				Benchmark: val.Benchmark,
				Samples:   val.Samples,
				Model:     val.Model,
			},
		}

		if value, ok := in[val.Model]; ok {
			if res.Cores == "" && value.Cores != "" {
				fmt.Println("Failed result cres", val.URL)
				continue Outer
			}
			for index := range res.Scores {
				if res.Scores[index] == "" && value.Scores[index] != "" {
					fmt.Println("Failed result scres", val.URL)
					continue Outer
				}
			}
			for index := range res.SegmentPerf {
				if res.SegmentPerf[index] == "" && value.SegmentPerf[index] != "" {
					fmt.Println("Failed result segments", val.URL)
					continue Outer
				}
			}
			for index := range res.SubResults {
				if res.SubResults[index] == "" && value.SubResults[index] != "" {
					fmt.Println("Failed result sub results", val.URL)
					continue Outer
				}
			}
		}

		// If all cecks succeed
		cpus[res.Model] = res

		out, err := json.MarshalIndent(cpus, "", "  ")
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
*/

/*
func getSSD(ctxt context.Context, cdp *chromedp.CDP, res *ssd, url string) {
	clor.Set(color.FgBlue)
	fmt.Println("Going to ", url)
	clor.Unset()
	if err := cp.Run(ctxt, chromedp.Navigate(url)); err != nil {
		clor.Set(color.FgRed)
		fmt.Println("Error navigating to ", url, err)
		clor.Unset()
		return
	}

	cRun(ctxt, chromedp.Sleep(time.Second*10))

	var wait syncWaitGroup
	wait.Add(3)

	go func() {
		defer wait.Done()
		cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
		defer cancel()
		clor.Set(color.FgCyan)
		fmt.Println("Trying cntroller")
		clor.Unset()
		if err := cp.Run(ctxt, chromedp.Text(`.cmp-cpt.medp.cmp-cpt-l`, &res.Controller, cdp.BySearch)); err != nil {
			clor.Set(color.BgRed)
			fmt.Println("Error getting cntroller for ", url, err)
			clor.Unset()
		} else {
			clor.Set(color.FgGreen)
			fmt.Println("Found cntroller")
			clor.Unset()
		}
	}()

	go func() {
		defer wait.Done()
		for i := 0; i < 9; i++ {
			cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
			defer cancel()
			clor.Set(color.FgCyan)
			fmt.Println(i, "Trying Subresult")
			clor.Unset()
			if err := cRun(ctxt, chromedp.Text(`.mcs-hl-col`, &res.SubResults[i], cdp.BySearch)); err != nil {
				clor.Set(color.BgHiRed)
				fmt.Print(i, "Error getting subresult", url, err)
				clor.Unset()
				fmt.Println()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found subresult")
				clor.Unset()
				func() {
					if err := cRun(ctxt, chromedp.SetAttributeValue(`.mcs-hl-col`, "class", "", cdp.BySearch)); err != nil {
						clor.Set(color.BgHiRed)
						fmt.Print(i, ".mc-hl-col", url, err)
						clor.Unset()
						fmt.Println()
					}
				}()
			}
		}
	}()

	go func() {
		defer wait.Done()
		for i := 0; i < 3; i++ {
			clor.Set(color.FgCyan)
			fmt.Println(i, "Trying Average")
			clor.Unset()
			err := func() error {
				cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
				defer cancel()
				if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3), &res.Averages[i])); err != nil {
					clor.Set(color.FgYellow)
					fmt.Println(i, "Error getting averages for ", url, err, "\ntrying yellow")
					clor.Unset()
					return err
				}
				return nil
			}()

			//if green fails, try yellow
			if err != nil {
				err = func() error {
					cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3), &res.Averages[i])); err != nil {
						clor.Set(color.FgRed)
						fmt.Println(i, "Error getting averages for ", url, err, "\ntrying red")
						clor.Unset()
						return err
					}
					return nil
				}()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found green")
				clor.Unset()
				cntinue
			}

			//if yellow fails, try red
			if err != nil {
				err = func() error {
					cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3), &res.Averages[i])); err != nil {
						clor.Set(color.BgRed)
						fmt.Print(i, "IMPORTANT!! Error getting averages for ", url, err)
						clor.Unset()
						fmt.Println()
						return err
					}
					return nil
				}()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found yellow")
				clor.Unset()
				cntinue
			}

			if err == nil {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found red")
				clor.Unset()
			}
		}
	}()

	wait.Wait()

	ssds = append(ssds, *res)
}

func getGPU(ctxt context.Context, c *chromedp.CDP, res *gpu, url string) {
	clor.Set(color.FgBlue)
	fmt.Println("Going to ", url)
	clor.Unset()
	if err := cRun(ctxt, chromedp.Navigate(url)); err != nil {
		clor.Set(color.FgRed)
		fmt.Println("Error navigating to ", url, err)
		clor.Unset()
		return
	}

	cRun(ctxt, chromedp.Sleep(time.Second*10))

	var wait syncWaitGroup
	wait.Add(2)
	go func() {
		for i := 0; i < 6; i++ {
			cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
			defer cancel()
			clor.Set(color.FgCyan)
			fmt.Println(i, "Trying Subresult")
			clor.Unset()
			if err := cRun(ctxt, chromedp.Text(`.mcs-hl-col`, &res.SubResults[i], cdp.BySearch)); err != nil {
				clor.Set(color.BgHiRed)
				fmt.Print(i, "Error getting subresult", url, err)
				clor.Unset()
				fmt.Println()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found subresult")
				clor.Unset()
				func() {
					if err := cRun(ctxt, chromedp.SetAttributeValue(`.mcs-hl-col`, "class", "", cdp.BySearch)); err != nil {
						clor.Set(color.BgHiRed)
						fmt.Print(i, ".mc-hl-col", url, err)
						clor.Unset()
						fmt.Println()
					}
				}()
			}
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 2; i++ {
			clor.Set(color.FgCyan)
			fmt.Println(i, "Trying Average")
			clor.Unset()
			err := func() error {
				cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
				defer cancel()
				if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3), &res.Averages[i])); err != nil {
					clor.Set(color.FgYellow)
					fmt.Println(i, "Error getting averages for ", url, err, "\ntrying yellow")
					clor.Unset()
					return err
				}
				return nil
			}()

			//if green fails, try yellow
			if err != nil {
				err = func() error {
					cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3), &res.Averages[i])); err != nil {
						clor.Set(color.FgRed)
						fmt.Println(i, "Error getting averages for ", url, err, "\ntrying red")
						clor.Unset()
						return err
					}
					return nil
				}()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found green")
				clor.Unset()
				cntinue
			}

			//if yellow fails, try red
			if err != nil {
				err = func() error {
					cxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := cRun(ctxt, chromedp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3), &res.Averages[i])); err != nil {
						clor.Set(color.BgRed)
						fmt.Print(i, "IMPORTANT!! Error getting averages for ", url, err)
						clor.Unset()
						fmt.Println()
						return err
					}
					return nil
				}()
			} else {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found yellow")
				clor.Unset()
				cntinue
			}

			if err == nil {
				clor.Set(color.FgGreen)
				fmt.Println(i, "found red")
				clor.Unset()
			}
		}
		wait.Done()
	}()

	wait.Wait()

	gpus = append(gpus, *res)
} */

/* func parseCSV(filename string) (out map[string]standard) {
	out = make(map[string]standard)
	copies := make(map[string]bool)

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	columns, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for i := 1; i < len(columns); i++ {
		if _, ok := copies[columns[i][3]]; ok {
			continue
		}
		s := make([]string, 8)
		for j := 1; j < len(columns[i]); j++ {
			s[j] = columns[i][j]
		}
		out[columns[i][3]] = standard{
			PartNum: s[1],
			Brand:   s[2],
			Model:   s[3],
			Rank: func() int {
				k, err := strconv.Atoi(s[4])
				if err != nil {
					fmt.Println(err, s[4])
				}
				return k
			}(),
			Benchmark: func() float32 {
				k, err := strconv.ParseFloat(s[5], 32)
				if err != nil {
					fmt.Println(err, s[5])
				}
				return float32(k)
			}(),
			Samples: func() int {
				k, err := strconv.Atoi(s[6])
				if err != nil {
					fmt.Println(err, s[6])
				}
				return k
			}(),
			URL: s[7],
		}
		copies[columns[i][3]] = true
	}

	return
}
 */