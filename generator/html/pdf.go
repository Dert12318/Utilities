package html

import (
	"context"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PdfOptional func(*page.PrintToPDFParams) *page.PrintToPDFParams

func (h *HTML) PDF(parentCtx context.Context, url, filepath, filename, html string, optional ...PdfOptional) ([]byte, error) {
	// log the CDP messages so that you can find the one to use.
	aCtx, aCancel := chromedp.NewRemoteAllocator(context.Background(), url)
	defer aCancel()

	ctx, cancel := chromedp.NewContext(aCtx)
	defer cancel()

	var resp []byte

	if err := chromedp.Run(ctx,
		// the navigation will trigger the "page.EventLoadEventFired" event too,
		// so we should add the listener after the navigation.
		chromedp.Navigate("about:blank"),
		// set the page content and wait until the page is loaded (including its resources).
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					// It's a good habit to remove the event listener if we don't need it anymore.
					cancel()
					wg.Done()
				}
			})

			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			build := page.PrintToPDF().
				WithDisplayHeaderFooter(true).
				WithFooterTemplate(`<div style="font-size:8px;width:100%;text-align:center;">(<span class="pageNumber"></span> / <span class="totalPages"></span>)</div>`).
				WithMarginBottom(0.4).
				WithMarginRight(0.4).
				WithMarginLeft(0.4).
				WithMarginTop(0.4).
				WithPrintBackground(false)

			for _, o := range optional {
				build = o(build)
			}

			buf, _, err := build.Do(ctx)
			if err != nil {
				return err
			}
			// return ioutil.WriteFile(fmt.Sprint(filepath+filename), buf, 0644)
			resp = buf
			return nil
		}),
	); err != nil {
		return resp, err
	}

	return resp, nil
}
