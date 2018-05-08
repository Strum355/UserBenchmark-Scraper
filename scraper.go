package main

import (
	"context"
	"fmt"
	"github.com/mafredri/cdp/protocol/runtime"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

var (
	conf = new(config)
	cpus = new(CPUs)
	gpus = new(GPUs)
	ssds = new(SSDs)
)

func init() {
	conf.loadConfig()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*50))
	defer cancel()

	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer domContent.Close()

	if err = c.Page.Enable(ctx); err != nil {
		fmt.Println(err)
		return
	}

	navArgs := page.NewNavigateArgs("http://www.userbenchmark.com/page/login")
	if _, err = c.Page.Navigate(ctx, navArgs); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = domContent.Recv(); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 2)

	/*   	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}  */

	fmt.Println("trying to insert username")

	if _, err = c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
		Expression: fmt.Sprintf(`document.querySelector('input[name="username"]').value = '%s'`, conf.User),
	}); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 2)

	if _, err = c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
		Expression: fmt.Sprintf(`document.querySelector('input[name="password"]').value = '%s'`, conf.Pass),
	}); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second)

	if _, err = c.Runtime.Evaluate(ctx, &runtime.EvaluateArgs{
		Expression: `document.querySelector('button[name="submit"]').click()`,
	}); err != nil {
		fmt.Println(err)
		return
	}

	/*
		// Get the outer HTML for the page.
		result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
			NodeID: &doc.Root.NodeID,
		})
		if err != nil {

			return
		}

		fmt.Printf("HTML: %s\n", result.OuterHTML) */
}
