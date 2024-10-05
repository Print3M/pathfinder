package crawler

import (
	"fmt"
	"pathfinder/src/store"
	"strings"
)

type Crawler struct {
	rootUrl      *store.Url
	noSubdomains bool
	noExternals  bool
	withAssets   bool
	headers      []string
}

func NewCrawler(rootUrl string, noSubdomains bool, noExternals bool, withAssets bool, headers []string) *Crawler {
	url, err := store.NewUrl(rootUrl)
	if err != nil {
		panic("Parse error")
	}

	return &Crawler{
		rootUrl:      url,
		noSubdomains: noSubdomains,
		noExternals:  noExternals,
		withAssets:   withAssets,
		headers:      headers,
	}
}

func (c *Crawler) RootUrl() *store.Url {
	return c.rootUrl
}

func (c *Crawler) ScrapUrlsFromUrl(collector *collector, targetUrl store.Url) []store.Url {
	rawPaths := collector.collectRawPaths(targetUrl.String())

	// Final filtering and processing
	var results []store.Url
	for _, rawPath := range rawPaths {
		rawPath = strings.TrimSpace(rawPath)

		if strings.HasPrefix(rawPath, "data:") {
			continue
		}

		url, err := c.rootUrl.Parse(rawPath)
		if err != nil {
			fmt.Printf("[-] Parse Error: %v\n", rawPath)
			continue
		}

		if c.noSubdomains && url.Host != c.rootUrl.Host {
			continue
		}

		if !strings.HasSuffix(url.Host, c.rootUrl.Host) {
			if c.noExternals {
				continue
			}

			url.IsExternal = true
		}

		if c.noExternals && !strings.HasSuffix(url.Host, c.rootUrl.Host) {
			continue
		}

		results = append(results, *url)
	}

	return results
}

func (c *Crawler) NewCollector() *collector {
	collector := newCollector(c.withAssets, c.headers)

	return collector
}
