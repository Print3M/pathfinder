package main

import (
	"fmt"
	"os"
	"pathfinder/src/cli"
	"pathfinder/src/crawler"
	"pathfinder/src/store"
	"pathfinder/src/workers"
	"time"
)

/*
	TODO:
	- add images
	- proxy
	- cookies
	- interrupt ctrl+c
	- rate limiting
*/

func crawlerWorker(c *crawler.Crawler, s *store.ScrapStore, input <-chan store.Url, done chan<- struct{}) {
	collector := crawler.NewCollector()

	// Work as long as input channel is open.
	for url := range input {
		urls := c.ScrapUrlsFromUrl(collector, url)

		for _, url := range urls {
			s.AddUrl(url)
		}

		done <- struct{}{}
	}
}

func runInitScrap(c *crawler.Crawler, s *store.ScrapStore) {
	s.AddUrl(*c.RootUrl())
	url, _ := s.GetNextUrlToVisit()
	collector := crawler.NewCollector()
	urls := c.ScrapUrlsFromUrl(collector, url)
	s.IncrementVisits()

	for _, v := range urls {
		s.AddUrl(v)
	}
}

func showStats(s *store.ScrapStore, start time.Time) {
	fmt.Println()
	fmt.Printf("Pages visited: %v\n", s.Visits())
	fmt.Printf("URLs  scraped: %v\n", s.CountScrapedUrls())
	fmt.Printf("Elapsed time : %v\n", time.Since(start))
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
				url, isFound := s.GetNextUrlToVisit()

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

func prepareOutputFile(name string) *os.File {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Open file error: %v\n", err)
		os.Exit(1)
	}

	return file
}

func start(flags *cli.Flags) {
	startTime := time.Now()
	c := crawler.NewCrawler(flags.Url, flags.NoSubdomains, flags.NoExternals, flags.WithAssets)

	// Prepare output file
	var file *os.File
	if len(flags.Output) > 0 {
		file = prepareOutputFile(flags.Output)
		defer file.Close()
	}
	s := store.NewStore(file, flags.Quiet)

	defer showStats(s, startTime)

	// First scrap to get initial input
	runInitScrap(c, s)
	if flags.NoRecursion {
		return
	}

	runCrawlerWorkers(c, s, flags.Threads)

}

func main() {
	cli.InitCli(start)
}
