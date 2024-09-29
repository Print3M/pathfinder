package main

import (
	"fmt"
	"scraper/cli"
	"scraper/crawler"
	"scraper/store"
	"time"

	"github.com/spf13/cobra"
)

/*
	TODO:
		- more tags than just a[href]
*/

type Worker struct {
	input  chan<- store.Url
	done   <-chan struct{}
	isIdle bool
}

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
	fmt.Printf("Scraped: %v\n", s.CountScraped())
}

func start() {
	c := crawler.New(flags.Url, flags.NoSubdomains)
	s := store.New()

	defer showStats(s)

	// Run init scrap
	initScrap(c, s)
	if flags.NoRecursion {
		return
	}

	// Run workers
	pool := make([]Worker, flags.Threads)
	for i := uint64(0); i < flags.Threads; i++ {
		done := make(chan struct{})
		input := make(chan store.Url)
		pool[i] = Worker{
			input:  input,
			done:   done,
			isIdle: true,
		}

		go crawlerWorker(c, s, input, done)
	}

	// Scheduler
	ticker := time.NewTicker(time.Millisecond * 50)
	idleCounter := len(pool)
	for {
		for i := 0; i < len(pool); i++ {
			if pool[i].isIdle {
				url, isFound := s.GetNextToScrap()

				if isFound {
					idleCounter--
					pool[i].isIdle = false
					pool[i].input <- url
					s.IncrementVisits()
				} else {
					if idleCounter == len(pool) {
						// There's no work to be done, exit.
						return
					}
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
