package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

const (
	domain  = "https://stackoverflow.com"
)

func main() {
	var nodes []*cdp.Node

	ctx := context.Background()
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("auto-open-devtools-for-tabs", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36"),
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, cc := chromedp.NewExecAllocator(ctx, options...)
	defer cc()
	// create context
	ctx, cancel := chromedp.NewContext(c)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(domain + "/questions/tagged/javascript%20reactjs?sort=Newest&filters=NoAcceptedAnswer&edited=true"),
		chromedp.Sleep(3 * time.Second),
		chromedp.Nodes(".question-hyperlink", &nodes, chromedp.NodeReady),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range nodes {
				url, _ := node.Attribute("href")
				if strings.HasPrefix(url, "https") == false {
					const selector = ".question p > a"
					var res []byte
					web := domain + url
					err := chromedp.Run(ctx,
						chromedp.Navigate(web),
						chromedp.Evaluate(fmt.Sprintf(`document.querySelectorAll("%s")`, selector), &res),
						chromedp.ActionFunc(func(ctx context.Context) error {
							// to avoid blocking behavior
							r := string(res)
							if r != "" && r != "{}" {
								var ids []cdp.NodeID
								err := chromedp.Run(ctx,
									chromedp.NodeIDs(selector, &ids),
									chromedp.ActionFunc(func(ctx context.Context) error {
										for _, id := range ids {
											attributes, err := dom.GetAttributes(id).Do(ctx)
											
											for i, attribute := range attributes {
												if attribute == "href" {
													a := attributes[i + 1]
													if strings.HasPrefix(a, "https://codesandbox.io") {
														fmt.Println(a)
														fmt.Println(web)
													}
												}

											}
											if err != nil {
												return err
											}
										}
										return nil
									}),
								)

								if err != nil {
									return err
								}
							}

							return nil
						}),
					)	
					if err != nil {
						return err
					}
				}
			}
			return nil
		}),
		chromedp.Sleep(1 * time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}
