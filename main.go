package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	domain  = "https://stackoverflow.com"
)

func main() {
	var nodes []*cdp.Node
	var buf []byte

	ctx := context.Background()
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
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
		chromedp.Sleep(6 * time.Second),
		chromedp.Nodes(".question-hyperlink", &nodes, chromedp.NodeReady),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range nodes {
				url, _ := node.Attribute("href")
				if strings.HasPrefix(url, "https") == false {
					fmt.Println(domain + url)
				}
			}
			fmt.Println(len(nodes))
			return nil
		}),
		chromedp.Sleep(1 * time.Second),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("stackoverflow.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
