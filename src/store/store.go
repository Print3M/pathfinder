package store

import (
	"fmt"
	"os"
	"sync"
)

type ScrapStore struct {
	toVisit map[string]Url
	scraped map[string]Url
	visits  uint
	file    *os.File

	quiet bool

	mu sync.RWMutex
}

func NewStore(file *os.File, quiet bool) *ScrapStore {
	return &ScrapStore{
		toVisit: make(map[string]Url),
		scraped: make(map[string]Url),
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
	s.mu.Lock()
	defer s.mu.Unlock()

	key := url.String()

	if _, ok := s.toVisit[key]; ok {
		return
	}

	if _, ok := s.scraped[key]; ok {
		return
	}

	if !s.quiet {
		fmt.Println(key)
	}

	s.toVisit[key] = url

	s.appendToFile(key)
}

func (s *ScrapStore) addExternalUrl(url Url) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := url.String()

	// TODO: Is checking required?
	if _, ok := s.scraped[key]; ok {
		return
	}

	if !s.quiet {
		fmt.Println(key)
	}

	// We don't visit external URL so add them directly to scraped URLs
	s.scraped[key] = url

	s.appendToFile(key)
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

	var key string
	var val Url
	for k, v := range s.toVisit {
		key = k
		val = v
		break
	}

	s.scraped[key] = val
	delete(s.toVisit, key)

	return val, true
}

func (s *ScrapStore) IncrementVisits() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visits++
}
