package main

import (
	"fmt"
	"scraper/cli"
	"scraper/crawler"
	"slices"
	"time"

	"github.com/spf13/cobra"
)

/*
	TODO:
		- remove GET params
		- more tags than just a[href]
*/

type Worker struct {
	input  chan<- string
	done   <-chan struct{}
	isIdle bool
}

var flags cli.Flags

func extractUrlsWorker(c *crawler.Crawler, input <-chan string, done chan<- struct{}) {
	for url := range input {
		// c.ChangeUrlStatus(url, crawler.Processed)
		urls := c.ScrapUrlsFromWeb(url)

		slices.Sort(urls)
		for _, url := range urls {
			c.AddUrl(url)
		}

		c.ChangeUrlStatus(url, crawler.Done)

		done <- struct{}{}
	}
}

func initScrap(c *crawler.Crawler) {
	targetUrl := c.RootUrl().String()
	c.AddUrl(targetUrl)
	urls := c.ScrapUrlsFromWeb(targetUrl)
	c.ChangeUrlStatus(targetUrl, crawler.Done)

	for _, url := range urls {
		c.AddUrl(url)
	}

	c.IncVisits()
}

func showStats(c *crawler.Crawler) {
	fmt.Printf("Visited: %v\n", c.Visits())
	fmt.Printf("Scraped: %v\n", len(c.Urls()))
}

func start() {
	c := crawler.New(flags.Url, flags.NoSubdomains)
	defer showStats(c)

	// Run init scrap
	initScrap(c)

	if flags.NoRecursion {
		return
	}

	// Run workers
	pool := make([]Worker, flags.Threads)
	for i := uint64(0); i < flags.Threads; i++ {
		done := make(chan struct{})
		input := make(chan string)
		pool[i] = Worker{
			input:  input,
			done:   done,
			isIdle: true,
		}

		go extractUrlsWorker(c, input, done)
	}

	// Scheduler
	ticker := time.NewTicker(time.Millisecond * 50)
	idleCounter := len(pool)
	visits := 0
	for {
		for i := 0; i < len(pool); i++ {
			if pool[i].isIdle {
				url, isFound := c.FindInitUrl()

				visits++
				fmt.Println(i, url, visits, idleCounter)

				if !isFound && idleCounter == len(pool) {
					// There's no work to be done, exit.
					return
				} else {
					c.ChangeUrlStatus(url.Url(), crawler.Processed)
					pool[i].isIdle = false
					idleCounter--
					pool[i].input <- url.Url()
				}
			} else {
				select {
				case <-pool[i].done:
					pool[i].isIdle = true
					idleCounter++
					break
				default:
					continue
				}
			}
		}

		<-ticker.C
	}
}

var rootCmd = &cobra.Command{
	Use:   "Urler",
	Short: "This is example short",
	Long:  "This is example long",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func main() {
	rootCmd.Flags().StringVarP(&flags.Url, "url", "u", "", "URL to scrap")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().Uint64VarP(&flags.Threads, "threads", "t", 10, "Number of concurrent threads")
	rootCmd.Flags().BoolVar(&flags.NoRecursion, "no-recursion", false, "Disable recursive scrapping")
	rootCmd.Flags().BoolVar(&flags.NoSubdomains, "no-subdomains", false, "Disable subdomain scrapping")
	rootCmd.Execute()
}
