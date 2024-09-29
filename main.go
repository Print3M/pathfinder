package main

import (
	"fmt"
	"scraper/cli"
	"scraper/crawler"
	"scraper/store"
	"scraper/workers"
	"time"

	"github.com/spf13/cobra"
)

/*
	TODO:
		- more tags than just a[href]
*/

var flags cli.Flags

func crawlerWorker(c *crawler.Crawler, s *store.ScrapStore, input <-chan store.Url, done chan<- struct{}) {
	for url := range input {
		urls := c.ScrapUrlsFromUrl(url)

		for _, url := range urls {
			s.AddToScrap(url)
		}

		done <- struct{}{}
	}
}

func initScrap(c *crawler.Crawler, s *store.ScrapStore) {
	s.AddToScrap(*c.RootUrl())
	url, _ := s.GetNextToScrap()
	urls := c.ScrapUrlsFromUrl(url)

	for _, v := range urls {
		s.AddToScrap(v)
	}
}

func showStats(s *store.ScrapStore) {
	fmt.Printf("Visited: %v\n", s.Visits())
	fmt.Printf("Scraped: %v\n", s.CountTotalStoredUrls())
}

func runCrawlerWorkers(c *crawler.Crawler, s *store.ScrapStore) {
	// Run workers
	pool := workers.NewPool(flags.Threads)
	pool.InitWorkers(func(input chan store.Url, done chan struct{}) {
		go crawlerWorker(c, s, input, done)
	})

	// Scheduler
	ticker := time.NewTicker(time.Millisecond * 50)
	for {
		for workerId := uint64(0); workerId < pool.Size; workerId++ {
			worker := pool.GetWorkerById(workerId)

			if worker.IsIdle() {
				url, isFound := s.GetNextToScrap()

				if isFound {
					worker.AssignJob(url)
					s.IncrementVisits()
				} else {
					if pool.AllWorkersIdle() {
						// There's no work to be done, exit.
						return
					}
				}
			} else {
				select {
				case <-worker.Done():
					worker.SetIdle()
					break
				default:
					continue
				}
			}
		}

		<-ticker.C
	}
}

func start() {
	c := crawler.New(flags.Url, flags.NoSubdomains)
	s := store.New()

	defer showStats(s)

	// First scrap to get initial input
	initScrap(c, s)
	if flags.NoRecursion {
		return
	}

	runCrawlerWorkers(c, s)
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
