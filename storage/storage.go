package storage

import (
	neturl "net/url"
	"sync"
)

type Url neturl.URL

func (u *Url) IsEqual(url neturl.URL) bool {
	return u.Scheme == url.Scheme && u.Host == url.Host && u.Path == url.Path
}

type ScrapStorage struct {
	toScrap []Url
	scraped []Url
	visits  uint64

	mu sync.RWMutex
}

func (s *ScrapStorage) AddToScrap(url neturl.URL) {
	s.mu.RLock()

	for _, item := range s.toScrap {
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
	defer s.mu.Unlock()

	s.toScrap = append(s.toScrap, Url(url))
}

func (s *ScrapStorage) GetNextToScrap() Url {
	s.mu.Lock()
	defer s.mu.Unlock()

	toScrap := s.toScrap[0]
	s.toScrap = s.toScrap[1:]
	s.scraped = append(s.scraped, toScrap)

	return toScrap
}

func (s *ScrapStorage) IncrementVisits() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visits++
}
