package crawler

import (
	"fmt"
	neturl "net/url"
	"pathfinder/src/store"
	"strings"

	"github.com/gocolly/colly"
)

type Crawler struct {
	rootUrl      neturl.URL
	noSubdomains bool
	withAssets   bool
}

func New(rootUrl string, noSubdomains bool, withAssets bool) *Crawler {
	url, err := neturl.Parse(rootUrl)
	if err != nil {
		panic("Parse error")
	}

	return &Crawler{
		rootUrl:      *url,
		noSubdomains: noSubdomains,
		withAssets:   withAssets,
	}
}

func (c *Crawler) RootUrl() *neturl.URL {
	return &c.rootUrl
}

func (c *Crawler) scrapRawUrls(targetUrl string) []string {
	// TODO: Optimize, spawn colly only once + mutex
	var urls []string
	collector := colly.NewCollector()

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	collector.OnHTML("form[action]", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("action"))
	})

	if c.withAssets {
		collector.OnHTML("img[src]", func(e *colly.HTMLElement) {
			urls = append(urls, e.Attr("src"))
		})

		collector.OnHTML("script[src]", func(e *colly.HTMLElement) {
			urls = append(urls, e.Attr("src"))
		})

		collector.OnHTML("link[href]", func(e *colly.HTMLElement) {
			urls = append(urls, e.Attr("href"))
		})
	}

	_ = collector.Visit(targetUrl)

	return urls
}

func (c *Crawler) ScrapUrlsFromUrl(targetUrl store.Url) []neturl.URL {
	rawUrls := c.scrapRawUrls(targetUrl.String())

	// Final filtering and processing
	var results []neturl.URL
	for _, rawUrl := range rawUrls {
		rawUrl = strings.TrimSpace(rawUrl)

		if strings.HasPrefix(rawUrl, "data:") {
			continue
		}

		url, err := c.rootUrl.Parse(rawUrl)
		if err != nil {
			fmt.Printf("[-] Parse Error: %v\n", rawUrl)
			continue
		}

		if c.noSubdomains && url.Host != c.rootUrl.Host {
			continue
		}

		if !strings.HasSuffix(url.Host, c.rootUrl.Host) {
			continue
		}

		results = append(results, *url)
	}

	return results
}
