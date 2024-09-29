package crawler

import (
	"fmt"
	neturl "net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type Status uint8

const (
	_ Status = iota
	Init
	Processed
	Done
)

type URL struct {
	url    string
	status Status
}

func (u URL) Url() string {
	return u.url
}

func (u URL) Status() Status {
	return u.status
}

type Options struct {
	noSubdomains bool
}

type Crawler struct {
	rootUrl neturl.URL
	options Options

	mu sync.RWMutex
}

func New(rootUrl string, noSubdomains bool) *Crawler {
	url, err := neturl.Parse(rootUrl)
	if err != nil {
		panic("Parse error")
	}

	return &Crawler{
		rootUrl: *url,
		options: Options{
			noSubdomains: noSubdomains,
		},
	}
}

func (c *Crawler) AddUrl(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.urls[url]; ok {
		return
	}

	// fmt.Println(url)
	c.urls[url] = URL{
		url:    url,
		status: Init,
	}
}

func (c *Crawler) Visits() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.visits
}

func (c *Crawler) Urls() map[string]URL {
	return c.urls
}

func (c *Crawler) RootUrl() *neturl.URL {
	return &c.rootUrl
}

func (c *Crawler) ChangeUrlStatus(url string, status Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.urls[url]

	if !ok {
		return
	}

	item.status = status
	c.urls[url] = item
}

func (c *Crawler) IncVisits() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.visits++
}

func (c *Crawler) ScrapUrlsFromWeb(targetWebUrl string) []string {
	collector := colly.NewCollector()
	var urls []string

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

		if c.options.noSubdomains {
			if href.Host != c.rootUrl.Host {
				return
			}
		} else {
			if !strings.HasSuffix(href.Host, c.rootUrl.Host) {
				return
			}
		}

		urls = append(urls, href.String())
	})

	// Scrap!
	err := collector.Visit(targetWebUrl)
	if err == nil {
		c.IncVisits()
	}

	return urls
}

func (c *Crawler) FindInitUrl() (URL, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for k, v := range c.urls {
		if v.status == Init {
			return c.urls[k], true
		}
	}

	return URL{}, false
}
