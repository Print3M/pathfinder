package workers

import (
	"pathfinder/src/crawler"
	"pathfinder/src/store"
	"time"
)

func crawlerWorker(c *crawler.Crawler, s *store.ScrapStore, input <-chan store.Url, done chan<- struct{}) {
	collector := c.NewCollector()

	// Work as long as input channel is open.
	for url := range input {
		urls := c.ScrapUrlsFromUrl(collector, url)

		for _, url := range urls {
			s.AddUrl(url)
		}

		done <- struct{}{}
	}
}

func RunCrawlerWorkers(c *crawler.Crawler, s *store.ScrapStore, threads uint, rate uint) {
	// Start workers
	pool := NewPool(threads)
	pool.InitWorkers(func(input chan store.Url, done chan struct{}) {
		go crawlerWorker(c, s, input, done)
	})

	// Prepare rate limiting
	var delay *time.Ticker
	if rate == 0 {
		delay = time.NewTicker(time.Nanosecond)
	} else {
		delay = time.NewTicker(time.Second / time.Duration(rate))
	}

	// Start scheduler
	ticker := time.NewTicker(time.Millisecond * 10)

Scheduler:
	for {
		for workerId := uint(0); workerId < pool.Size; workerId++ {
			worker := pool.GetWorkerById(workerId)
			worker.Update()

			if !worker.IsIdle() {
				continue
			}

			url, isFound := s.GetNextUrlToVisit()

			if !isFound {
				if pool.AreAllWorkersIdle() {
					// There's no work to be done, exit.
					pool.ShutdownAllWorkers()

					break Scheduler
				}

				continue
			}

			<-delay.C
			worker.AssignJob(url)
			s.IncrementVisits()
		}

		<-ticker.C
	}
}
