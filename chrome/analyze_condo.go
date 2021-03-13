package chrome

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

func loopAnalyzeCondo(ctx context.Context, condoUrls []string, maxConcurrent int) []*Condo {
	concurrentGoroutines := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	condos := make([]*Condo, 0)

	for _, url := range condoUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			condo := analyzeCondo(ctx, url)
			if condo != nil {
				condos = append(condos, condo)
			}
			<-concurrentGoroutines
		}(url)
	}

	wg.Wait()

	return condos
}

func analyzeCondo(ctx context.Context, condoUrl string) *Condo {
	if !isCondoUrl(condoUrl) {
		return nil
	}

	log.Infof("Now is %s\n", condoUrl)

	c := &Condo{Url: condoUrl}
	facilities := make([]string, 0)
	var leftNodes []*cdp.Node
	var rightNodes []*cdp.Node

	ctxNew, _ := chromedp.NewContext(ctx)

	err := chromedp.Run(ctxNew,
		chromedp.Navigate(condoUrl),
		// Name
		chromedp.Text(`div.propertytitlecontainer > div > h1`, &c.Name),
		// address
		chromedp.Text(`div.propertytitlecontainer > div.propcol1 > :nth-child(2)`, &c.Address),
		chromedp.Nodes(`div#tabs-property-facilities > :nth-child(2) > ul > li`, &leftNodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.Nodes(`div#tabs-property-facilities > :nth-child(3) > ul > li`, &rightNodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var fac string
			var facStr string
			for i := 1; i <= len(leftNodes); i++ {
				sel := fmt.Sprintf("div#tabs-property-facilities > :nth-child(2) > ul > :nth-child(%d)", i)
				err := chromedp.Text(sel, &fac).Do(ctx)
				if err != nil {
					return err
				}
				facilities = append(facilities, fac)
				facStr += fmt.Sprintf("| %s", fac)
			}

			for i := 1; i < len(rightNodes); i++ {
				sel := fmt.Sprintf("div#tabs-property-facilities > :nth-child(3) > ul > :nth-child(%d)", i)
				err := chromedp.Text(sel, &fac).Do(ctx)
				if err != nil {
					return err
				}
				facilities = append(facilities, fac)
				facStr += fmt.Sprintf("| %s", fac)
			}

			c.FacString = facStr

			return nil
		}),
	)
	if err != nil {
		log.Errorf("analyze condo %v", err)
		return nil
	}

	c.Facility = analyzeFac(facilities)

	return c
}
