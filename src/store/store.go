package store

import (
	"fmt"
	"os"
	"sync"
)

type ScrapStore struct {
	toVisit []Url
	scraped []Url
	visits  uint
	file    *os.File

	quiet bool

	mu sync.RWMutex
}

func NewStore(file *os.File, quiet bool) *ScrapStore {
	return &ScrapStore{
		toVisit: make([]Url, 0, 4096),
		scraped: make([]Url, 0, 4096),
		file:    file,
		quiet:   quiet,
	}
}

func (s *ScrapStore) appendToFile(v string) {
	if s.file != nil {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.file.Write([]byte(v + "\n"))
	}
}

func (s *ScrapStore) Visits() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.visits
}

func (c *ScrapStore) AddUrl(url Url) {
	if url.IsExternal {
		c.addExternalUrl(url)
	} else {
		c.addUrlToVisit(url)
	}
}

func (s *ScrapStore) addUrlToVisit(url Url) {
	s.mu.RLock()

	for _, item := range s.toVisit {
		if item.IsEqual(url) {
			s.mu.RUnlock()
			return
		}
	}

	for _, item := range s.scraped {
		if item.IsEqual(url) {
			s.mu.RUnlock()
			return
		}
	}

	s.mu.RUnlock()
	s.mu.Lock()

	if !s.quiet {
		fmt.Println(url.String())
	}

	s.toVisit = append(s.toVisit, url)

	s.mu.Unlock()

	s.appendToFile(url.String())
}

func (s *ScrapStore) addExternalUrl(url Url) {
	s.mu.RLock()

	for _, item := range s.scraped {
		if item.IsEqual(url) {
			s.mu.RUnlock()
			return
		}
	}

	s.mu.RUnlock()
	s.mu.Lock()

	if !s.quiet {
		fmt.Println(url.String())
	}

	// We don't visit external URL so add them directly to scraped URLs
	s.scraped = append(s.scraped, url)

	s.mu.Unlock()

	s.appendToFile(url.String())
}

func (s *ScrapStore) CountScrapedUrls() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint(len(s.scraped))
}

func (s *ScrapStore) GetNextUrlToVisit() (Url, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.toVisit) == 0 {
		return Url{}, false
	}

	toScrap := s.toVisit[0]
	s.toVisit = s.toVisit[1:]
	s.scraped = append(s.scraped, toScrap)

	return toScrap, true
}

func (s *ScrapStore) IncrementVisits() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visits++
}
