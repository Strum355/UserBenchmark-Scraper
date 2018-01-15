package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	cdp "github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"github.com/fatih/color"
)

var (
	conf = new(config)
)

func main() {
	conf.loadConfig()

	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := cdp.New(ctxt, cdp.WithRunnerOptions(
		runner.HeadlessPathPort("/usr/bin/google-chrome-stable", 9222), 
		runner.Flag("headless", false)), 
		cdp.WithErrorf(func(s string, v ...interface{}) {
			if strings.Contains(s, "could not perform") || strings.Contains(s, "could not get"){
				return
			}
			log.Printf("error: "+s, v...)
		}))
	if err != nil {	
		fmt.Println(err)
		return
	}

	start(ctxt, c)

	if err := c.Shutdown(ctxt); err != nil {
		fmt.Println(err)
	}
}

/*
	cdp.OuterHTML(`html`, &html, cdp.ByQuery),
*/

func start(ctx context.Context, c *cdp.CDP) {
	fmt.Println("STARTING")
	if err := login(ctx, c); err != nil {
		fmt.Println("Couldnt login", err)
		return
	}

	var html string
	c.Run(ctx, cdp.Tasks{
		cdp.Navigate(`http://gpu.userbenchmark.com/Nvidia-GTX-1070/Rating/3609`),
		cdp.OuterHTML(`html`, &html, cdp.ByQuery),
	})

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(`.mcs-hl-col`).Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
	//doCPU(ctx)

	/* 	time.Sleep(time.Minute*30)

	   	for name, val := range parseCSV("CPU_UserBenchmarks.csv") {
	   		var res cpu
	   		res.Name = name
	   		getCPU(ctx, &res, val)
	   		out, err := json.MarshalIndent(cpus, "", "  ")
	   		if err != nil {
	   			fmt.Println("JSON: ", err)
	   			return
	   		}

	   		err = ioutil.WriteFile("CPU_DATA.json", out, 0644)
	   		if err != nil {
	   			fmt.Println("WRITE: ", err)
	   		}
	   	} */

	/* 	time.Sleep(time.Minute*30)

	   	for name, val := range parseCSV("GPU_UserBenchmarks.csv") {
	   		var res gpu
	   		res.Name = name
	   		getGPU(ctx, &res, val)
	   		out, err := json.MarshalIndent(gpus, "", "  ")
	   		if err != nil {
	   			fmt.Println("JSON: ", err)
	   			return
	   		}

	   		err = ioutil.WriteFile("GPU_DATA.json", out, 0644)
	   		if err != nil {
	   			fmt.Println("WRITE: ", err)
	   		}
	   	}*/
}

func login(ctx context.Context, c *cdp.CDP) error {
	return c.Run(ctx, cdp.Tasks{
		cdp.Navigate(`http://www.userbenchmark.com/page/login`),
		cdp.SetValue(`input[name="username"]`, conf.User),
		cdp.Sleep(time.Second * 1),
		cdp.SetValue(`input[name="password"]`, conf.Pass),
		cdp.Sleep(time.Second * 1),
		cdp.Click(`button[name="submit"]`),
		cdp.Sleep(time.Second * 2),
	})
}

func doCPU(ctx context.Context, c *cdp.CDP) {
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
		getCPU(ctx, c, &res, val.URL)

		if value, ok := in[val.Model]; ok {
			if res.Cores == "" && value.Cores != "" {
				fmt.Println("Failed result cores", val.URL)
				continue Outer
			}
			for index := range res.Scores {
				if res.Scores[index] == "" && value.Scores[index] != "" {
					fmt.Println("Failed result scores", val.URL)
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

		// If all checks succeed
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

func getCPU(ctxt context.Context, c *cdp.CDP, res *cpu, url string) {
	color.Set(color.FgBlue)
	fmt.Println("Going to ", url)
	color.Unset()
	if err := c.Run(ctxt, cdp.Navigate(url)); err != nil {
		color.Set(color.FgRed)
		fmt.Println("Error navigating to ", url, err)
		color.Unset()
	}

	c.Run(ctxt, cdp.Sleep(time.Second*3))

	var wait sync.WaitGroup
	wait.Add(3)

	go func() {
		defer wait.Done()
		ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
		defer cancel()
		color.Set(color.FgCyan)
		fmt.Println("Trying cores")
		color.Unset()
		if err := c.Run(ctxt, cdp.Text(`.cmp-cpt.tallp.cmp-cpt-l`, &res.Cores, cdp.ByQuery)); err != nil {
			color.Set(color.BgRed)
			fmt.Println("Error getting cores for ", url, err)
			color.Unset()
		} else {
			color.Set(color.FgGreen)
			fmt.Println("Found cores")
			color.Unset()
		}
	}()

	go func() {
		defer wait.Done()
		var waitIn sync.WaitGroup
		for i := 0; i < 3; i++ {
			waitIn.Add(2)
			go func(i int) {
				defer waitIn.Done()
				color.Set(color.FgCyan)
				fmt.Println(i, "Trying green scores")
				color.Unset()
				err := func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3), &res.Scores[i], cdp.ByQuery)); err != nil {
						color.Set(color.FgYellow)
						fmt.Println(i, "Error getting scores for ", url, err, "\ntrying yellow")
						color.Unset()
						return err
					}
					return nil
				}()

				//if green fails, try yellow
				if err != nil {
					err = func() error {
						ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
						defer cancel()
						if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3), &res.Scores[i], cdp.ByQuery)); err != nil {
							color.Set(color.FgRed)
							fmt.Println(i, "Error getting scores for ", url, err, "\ntrying red")
							color.Unset()
							return err
						}
						return nil
					}()
				} else {
					color.Set(color.FgGreen)
					fmt.Println(i, "found green")
					color.Unset()
					return
				}

				//if yellow fails, try red
				if err != nil {
					err = func() error {
						ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
						defer cancel()
						if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3), &res.Scores[i], cdp.ByQuery)); err != nil {
							color.Set(color.BgRed)
							fmt.Println(i, "IMPORTANT!! Error getting scores for ", url, err)
							color.Unset()
							return err
						}
						return nil
					}()
				} else {
					color.Set(color.FgGreen)
					fmt.Println(i, "found yellow")
					color.Unset()
					return
				}

				if err == nil {
					color.Set(color.FgGreen)
					fmt.Println(i, "found red")
					color.Unset()
				}
			}(i)

			go func(i int) {
				defer waitIn.Done()
				color.Set(color.FgCyan)
				fmt.Println(i, "Trying performance")
				color.Unset()
				ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
				defer cancel()
				if err := c.Run(ctxt, cdp.Text(`.bsc-w.text-left.semi-strong`, &res.SegmentPerf[i], cdp.ByQuery)); err != nil {
					color.Set(color.FgHiRed)
					fmt.Println(i, "Error getting performance for ", url, err)
					color.Unset()
				} else {
					color.Set(color.FgGreen)
					fmt.Println(i, "found performance")
					color.Unset()
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*10)
					defer cancel()
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.bsc-w.text-left.semi-strong`, "class", "", cdp.ByQuery)); err != nil {
						color.Set(color.FgHiRed)
						fmt.Println(i, `.bsc-w.text-left.semi-strong`, url, err)
						color.Unset()
					}
				}
			}(i)
			waitIn.Wait()
		}
	}()

	go func() {
		for i := 0; i < 9; i++ {
			ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
			defer cancel()
			color.Set(color.FgCyan)
			fmt.Println(i, "Trying Subresult")
			color.Unset()
			if err := c.Run(ctxt, cdp.Text(`.mcs-hl-col`, &res.SubResults[i], cdp.ByQuery)); err != nil {
				color.Set(color.BgHiRed)
				fmt.Print(i, "Error getting subresult", url, err)
				color.Unset()
				fmt.Println()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found subresult")
				color.Unset()
				func() {
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-hl-col`, "class", "", cdp.ByQuery)); err != nil {
						color.Set(color.BgHiRed)
						fmt.Print(i, ".mcs-hl-col", url, err)
						color.Unset()
						fmt.Println()
					}
				}()
			}
		}
		wait.Done()
	}()

	wait.Wait()

	for i, val := range res.SegmentPerf {
		res.SegmentPerf[i] = strings.Trim(strings.Replace(strings.Replace(val, "\t", "", -1), "\n", " ", -1), " ")
	}
}

func getSSD(ctxt context.Context, c *cdp.CDP, res *ssd, url string) {
	color.Set(color.FgBlue)
	fmt.Println("Going to ", url)
	color.Unset()
	if err := c.Run(ctxt, cdp.Navigate(url)); err != nil {
		color.Set(color.FgRed)
		fmt.Println("Error navigating to ", url, err)
		color.Unset()
		return
	}

	c.Run(ctxt, cdp.Sleep(time.Second*10))

	var wait sync.WaitGroup
	wait.Add(3)

	go func() {
		defer wait.Done()
		ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
		defer cancel()
		color.Set(color.FgCyan)
		fmt.Println("Trying controller")
		color.Unset()
		if err := c.Run(ctxt, cdp.Text(`.cmp-cpt.medp.cmp-cpt-l`, &res.Controller, cdp.ByQuery)); err != nil {
			color.Set(color.BgRed)
			fmt.Println("Error getting controller for ", url, err)
			color.Unset()
		} else {
			color.Set(color.FgGreen)
			fmt.Println("Found controller")
			color.Unset()
		}
	}()

	go func() {
		defer wait.Done()
		for i := 0; i < 9; i++ {
			ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
			defer cancel()
			color.Set(color.FgCyan)
			fmt.Println(i, "Trying Subresult")
			color.Unset()
			if err := c.Run(ctxt, cdp.Text(`.mcs-hl-col`, &res.SubResults[i], cdp.ByQuery)); err != nil {
				color.Set(color.BgHiRed)
				fmt.Print(i, "Error getting subresult", url, err)
				color.Unset()
				fmt.Println()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found subresult")
				color.Unset()
				func() {
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-hl-col`, "class", "", cdp.ByQuery)); err != nil {
						color.Set(color.BgHiRed)
						fmt.Print(i, ".mcs-hl-col", url, err)
						color.Unset()
						fmt.Println()
					}
				}()
			}
		}
	}()

	go func() {
		defer wait.Done()
		for i := 0; i < 3; i++ {
			color.Set(color.FgCyan)
			fmt.Println(i, "Trying Average")
			color.Unset()
			err := func() error {
				ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
				defer cancel()
				if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3), &res.Averages[i])); err != nil {
					color.Set(color.FgYellow)
					fmt.Println(i, "Error getting averages for ", url, err, "\ntrying yellow")
					color.Unset()
					return err
				}
				return nil
			}()

			//if green fails, try yellow
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3), &res.Averages[i])); err != nil {
						color.Set(color.FgRed)
						fmt.Println(i, "Error getting averages for ", url, err, "\ntrying red")
						color.Unset()
						return err
					}
					return nil
				}()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found green")
				color.Unset()
				continue
			}

			//if yellow fails, try red
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3), &res.Averages[i])); err != nil {
						color.Set(color.BgRed)
						fmt.Print(i, "IMPORTANT!! Error getting averages for ", url, err)
						color.Unset()
						fmt.Println()
						return err
					}
					return nil
				}()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found yellow")
				color.Unset()
				continue
			}

			if err == nil {
				color.Set(color.FgGreen)
				fmt.Println(i, "found red")
				color.Unset()
			}
		}
	}()

	wait.Wait()

	ssds = append(ssds, *res)
}

func getGPU(ctxt context.Context, c *cdp.CDP, res *gpu, url string) {
	color.Set(color.FgBlue)
	fmt.Println("Going to ", url)
	color.Unset()
	if err := c.Run(ctxt, cdp.Navigate(url)); err != nil {
		color.Set(color.FgRed)
		fmt.Println("Error navigating to ", url, err)
		color.Unset()
		return
	}

	c.Run(ctxt, cdp.Sleep(time.Second*10))

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		for i := 0; i < 6; i++ {
			ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
			defer cancel()
			color.Set(color.FgCyan)
			fmt.Println(i, "Trying Subresult")
			color.Unset()
			if err := c.Run(ctxt, cdp.Text(`.mcs-hl-col`, &res.SubResults[i], cdp.ByQuery)); err != nil {
				color.Set(color.BgHiRed)
				fmt.Print(i, "Error getting subresult", url, err)
				color.Unset()
				fmt.Println()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found subresult")
				color.Unset()
				func() {
					if err := c.Run(ctxt, cdp.SetAttributeValue(`.mcs-hl-col`, "class", "", cdp.ByQuery)); err != nil {
						color.Set(color.BgHiRed)
						fmt.Print(i, ".mcs-hl-col", url, err)
						color.Unset()
						fmt.Println()
					}
				}()
			}
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 2; i++ {
			color.Set(color.FgCyan)
			fmt.Println(i, "Trying Average")
			color.Unset()
			err := func() error {
				ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
				defer cancel()
				if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pgbg`, i+3), &res.Averages[i])); err != nil {
					color.Set(color.FgYellow)
					fmt.Println(i, "Error getting averages for ", url, err, "\ntrying yellow")
					color.Unset()
					return err
				}
				return nil
			}()

			//if green fails, try yellow
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.pybg`, i+3), &res.Averages[i])); err != nil {
						color.Set(color.FgRed)
						fmt.Println(i, "Error getting averages for ", url, err, "\ntrying red")
						color.Unset()
						return err
					}
					return nil
				}()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found green")
				color.Unset()
				continue
			}

			//if yellow fails, try red
			if err != nil {
				err = func() error {
					ctxt, cancel := context.WithTimeout(ctxt, time.Second*20)
					defer cancel()
					if err := c.Run(ctxt, cdp.Text(fmt.Sprintf(`.para-m-t.uc-table.table-no-border > thead > tr > td:nth-child(%d) .mcs-caption.prbg`, i+3), &res.Averages[i])); err != nil {
						color.Set(color.BgRed)
						fmt.Print(i, "IMPORTANT!! Error getting averages for ", url, err)
						color.Unset()
						fmt.Println()
						return err
					}
					return nil
				}()
			} else {
				color.Set(color.FgGreen)
				fmt.Println(i, "found yellow")
				color.Unset()
				continue
			}

			if err == nil {
				color.Set(color.FgGreen)
				fmt.Println(i, "found red")
				color.Unset()
			}
		}
		wait.Done()
	}()

	wait.Wait()

	gpus = append(gpus, *res)
}

func parseCSV(filename string) (out map[string]standard) {
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

func isIn(name string) bool {
	for _, r := range gpus {
		if r.Name == name {
			return true
		}
	}
	return false
}
