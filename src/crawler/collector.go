package crawler

import (
	"math/rand/v2"
	"strings"

	"github.com/gocolly/colly"
)

type collector struct {
	engine     *colly.Collector
	headers    []string
	withAssets bool
}

func newCollector(withAssets bool, headers []string) *collector {
	engine := colly.NewCollector()

	c := collector{
		engine:     engine,
		withAssets: withAssets,
		headers:    headers,
	}

	// Prepare extra HTTP headers
	if len(headers) > 0 {
		c.engine.OnRequest(func(r *colly.Request) {
			for _, header := range headers {
				parts := strings.Split(header, ":")

				if len(parts) < 2 {
					continue
				}

				value := strings.Join(parts[1:], ":")
				r.Headers.Set(parts[0], value)
			}
		})
	}

	return &c
}

func (c *collector) addOnHtml(paths *[]string, querySelector string, attrName string) {
	c.engine.OnHTML(querySelector, func(e *colly.HTMLElement) {
		*paths = append(*paths, e.Attr(attrName))
	})
}

func (c *collector) getRandomUserAgent() string {
	return userAgents[rand.IntN(len(userAgents))]
}

func (c *collector) collectRawPaths(targetUrl string) []string {
	var paths []string
	c.addOnHtml(&paths, "a[href]", "href")
	c.addOnHtml(&paths, "form[action]", "action")
	c.addOnHtml(&paths, "iframe[src]", "src")
	c.addOnHtml(&paths, "area[href]", "href")
	c.addOnHtml(&paths, "base[href]", "href")

	if c.withAssets {
		c.addOnHtml(&paths, "area[href]", "href")
		c.addOnHtml(&paths, "img[src]", "src")
		c.addOnHtml(&paths, "script[src]", "src")
		c.addOnHtml(&paths, "link[href]", "href")
		c.addOnHtml(&paths, "embed[src]", "src")
		c.addOnHtml(&paths, "audio[src]", "src")
		c.addOnHtml(&paths, "object[data]", "data")
		c.addOnHtml(&paths, "video[src]", "src")
		c.addOnHtml(&paths, "track[src]", "src")
		c.addOnHtml(&paths, "source[src]", "src")
	}

	c.engine.UserAgent = c.getRandomUserAgent()

	_ = c.engine.Visit(targetUrl)

	return paths
}

var userAgents = [...]string{
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; rv:40.0) Gecko/20100101 Firefox/40.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.91 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.0; rv:38.0) Gecko/20100101 Firefox/38.0",
	"Mozilla/5.0 (Windows NT 6.1; Trident/7.0; FunWebProducts; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; rv:29.0) Gecko/20100101 Firefox/29.0",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0; EIE10;ENUSWOL)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:38.0) Gecko/20100101 Firefox/38.0 SeaMonkey/2.35",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Win64; x64; Trident/7.0)",
	"Mozilla/5.0 (Windows NT 5.1; rv:32.0) Gecko/20100101 Firefox/32.0",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
}
