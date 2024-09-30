package main

import (
	"fmt"
	"scraper/src/cli"
	"scraper/src/crawler"
	"scraper/src/store"
	"scraper/src/workers"
	"time"
)

func crawlerWorker(c *crawler.Crawler, s *store.ScrapStore, input <-chan store.Url, done chan<- struct{}) {
	// Work as long as input channel is open.
	for url := range input {
		urls := c.ScrapUrlsFromUrl(url)

		for _, url := range urls {
			s.AddToScrap(url)
		}

		done <- struct{}{}
	}
}

func runInitScrap(c *crawler.Crawler, s *store.ScrapStore) {
	s.AddToScrap(*c.RootUrl())
	url, _ := s.GetNextToScrap()
	urls := c.ScrapUrlsFromUrl(url)
	s.IncrementVisits()

	for _, v := range urls {
		s.AddToScrap(v)
	}
}

func showStats(s *store.ScrapStore) {
	fmt.Printf("Visited: %v\n", s.Visits())
	fmt.Printf("Scraped: %v\n", s.CountTotalStoredUrls())
}

func runCrawlerWorkers(c *crawler.Crawler, s *store.ScrapStore, threads uint64) {
	// Start workers
	pool := workers.NewPool(threads)
	pool.InitWorkers(func(input chan store.Url, done chan struct{}) {
		go crawlerWorker(c, s, input, done)
	})

	// Start scheduler
	ticker := time.NewTicker(time.Millisecond * 50)

Scheduler:
	for {
		for workerId := uint64(0); workerId < pool.Size; workerId++ {
			worker := pool.GetWorkerById(workerId)
			worker.Update()

			if worker.IsIdle() {
				url, isFound := s.GetNextToScrap()

				if isFound {
					worker.AssignJob(url)
					s.IncrementVisits()
				} else {
					if pool.AreAllWorkersIdle() {
						// There's no work to be done, exit.
						pool.ShutdownAllWorkers()

						break Scheduler
					}
				}
			}
		}

		<-ticker.C
	}
}

func start(flags *cli.Flags) {
	c := crawler.New(flags.Url, flags.NoSubdomains, flags.WithAssets)
	s := store.New()

	defer showStats(s)

	// First scrap to get initial input
	runInitScrap(c, s)
	if flags.NoRecursion {
		return
	}

	runCrawlerWorkers(c, s, flags.Threads)

	fmt.Printf("Exit....")
}

func main() {
	cli.InitCli(start)
}
