package crawler

import (
	"fmt"
	neturl "net/url"
	"scraper/store"
	"strings"

	"github.com/gocolly/colly"
)

type Crawler struct {
	rootUrl      neturl.URL
	noSubdomains bool
}

func New(rootUrl string, noSubdomains bool) *Crawler {
	url, err := neturl.Parse(rootUrl)
	if err != nil {
		panic("Parse error")
	}

	return &Crawler{
		rootUrl:      *url,
		noSubdomains: noSubdomains,
	}
}

func (c *Crawler) RootUrl() *neturl.URL {
	return &c.rootUrl
}

func (c *Crawler) ScrapUrlsFromUrl(targetUrl store.Url) []neturl.URL {
	collector := colly.NewCollector()
	var urls []neturl.URL

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		rawHref := e.Attr("href")
		rawHref = strings.TrimSpace(rawHref)

		href, err := c.rootUrl.Parse(rawHref)
		if err != nil {
			fmt.Printf("[-] Parse Error: %v\n", rawHref)
			return
		}

		// Don't check fragments
		href.Fragment = ""

		if c.noSubdomains {
			if href.Host != c.rootUrl.Host {
				return
			}
		} else {
			if !strings.HasSuffix(href.Host, c.rootUrl.Host) {
				return
			}
		}

		urls = append(urls, *href)
	})

	// Scrap!
	_ = collector.Visit(targetUrl.String())

	return urls
}
