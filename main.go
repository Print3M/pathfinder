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
	- proxy
	- interrupt ctrl+c
	- change slice to map
*/

func runInitScrap(c *crawler.Crawler, s *store.ScrapStore) {
	s.AddUrl(*c.RootUrl())
	url, _ := s.GetNextUrlToVisit()
	collector := c.NewCollector()
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

func start(flags *cli.Flags) {
	startTime := time.Now()
	c := crawler.NewCrawler(flags.Url, flags.NoSubdomains, flags.NoExternals, flags.WithAssets, flags.Headers)

	// Prepare output file
	var file *os.File
	if len(flags.Output) > 0 {
		file = cli.PrepareOutputFile(flags.Output)
		defer file.Close()
	}
	s := store.NewStore(file, flags.Quiet)

	defer showStats(s, startTime)

	// First scrap to get initial input
	runInitScrap(c, s)
	if flags.NoRecursion {
		return
	}

	workers.RunCrawlerWorkers(c, s, flags.Threads, flags.Rate)
}

func main() {
	cli.InitCli(start)
}
