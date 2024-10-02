package store

import (
	"fmt"
	"os"
	"sync"
)

type ScrapStore struct {
	toVisit []Url
	scraped []Url
	visits  uint64
	file    *os.File

	quiet bool

	mu sync.RWMutex
}

func NewStore(file *os.File, quiet bool) *ScrapStore {
	return &ScrapStore{
		toVisit: make([]Url, 0),
		scraped: make([]Url, 0),
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

func (s *ScrapStore) Visits() uint64 {
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

	parsed := Url(url)
	s.toVisit = append(s.toVisit, parsed)

	s.mu.Unlock()

	s.appendToFile(parsed.String())
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
	parsed := Url(url)
	s.scraped = append(s.scraped, parsed)

	s.mu.Unlock()

	s.appendToFile(parsed.String())
}

func (s *ScrapStore) CountScrapedUrls() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint64(len(s.scraped))
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
