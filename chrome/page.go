package chrome

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
)

func Navigate(ctx context.Context, c *cdp.Client, pageURL string) error {
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContent.Close()

	if err = c.Page.Enable(ctx); err != nil {
		return err
	}

	if _, err = c.Page.Navigate(ctx,
		page.NewNavigateArgs(pageURL)); err != nil {
		return err
	}

	if _, err = domContent.Recv(); err != nil {
		return err
	}

	return nil
}
