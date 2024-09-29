package store

import (
	"fmt"
	neturl "net/url"
	"sync"
)

type Url neturl.URL

func (u *Url) IsEqual(url neturl.URL) bool {
	return u.Scheme == url.Scheme && u.Host == url.Host && u.Path == url.Path
}

func (u *Url) String() string {
	return u.Scheme + "://" + u.Host + u.Path
}

type ScrapStore struct {
	toScrap []Url
	scraped []Url
	visits  uint64

	mu sync.RWMutex
}

func New() *ScrapStore {
	return &ScrapStore{
		toScrap: make([]Url, 0),
		scraped: make([]Url, 0),
	}
}

func (s *ScrapStore) Visits() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.visits
}

func (s *ScrapStore) CountTotalStoredUrls() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint64(len(s.scraped) + len(s.toScrap))
}

func (s *ScrapStore) CountToScrap() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return uint64(len(s.toScrap))
}

func (s *ScrapStore) AddToScrap(url neturl.URL) {
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

	fmt.Println(url.String())
	s.toScrap = append(s.toScrap, Url(url))
}

func (s *ScrapStore) GetNextToScrap() (Url, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.toScrap) == 0 {
		return Url{}, false
	}

	toScrap := s.toScrap[0]
	s.toScrap = s.toScrap[1:]
	s.scraped = append(s.scraped, toScrap)

	return toScrap, true
}

func (s *ScrapStore) IncrementVisits() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visits++
}
