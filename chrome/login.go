package chrome

import (
	"context"
	"fmt"
	"time"

	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/Strum355/UserBenchmark-Scraper/config"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
)

// Start conects to the headless chrome at URL
func Start(ctx context.Context, URL string) (c *cdp.Client, conn *rpcc.Conn, err error) {
	devt := devtool.New(URL)

	var pt *devtool.Target
	pt, err = devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return
		}
	}

	conn, err = rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return
	}

	c = cdp.NewClient(conn)
	return
}

// Login simply logs in I guess
func Login(ctx context.Context, c *cdp.Client, conf *config.Config) (err error) {
	if err = Navigate(ctx, c, "http://www.userbenchmark.com/page/login"); err != nil {
		return
	}

	args := runtime.NewEvaluateArgs(
		fmt.Sprintf(`document.querySelector('input[name="username"]').value = '%s'`, conf.User),
	)
	if _, err = c.Runtime.Evaluate(ctx, args); err != nil {
		return
	}

	time.Sleep(time.Second * 2)

	args = runtime.NewEvaluateArgs(
		fmt.Sprintf(`document.querySelector('input[name="password"]').value = '%s'`, conf.Pass),
	)
	if _, err = c.Runtime.Evaluate(ctx, args); err != nil {
		return
	}

	time.Sleep(time.Second)

	args = runtime.NewEvaluateArgs(`document.querySelector('button[name="submit"]').click()`)
	if _, err = c.Runtime.Evaluate(ctx, args); err != nil {
		return
	}

	time.Sleep(time.Second * 2)

	return
}
