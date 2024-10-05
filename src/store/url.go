package store

import (
	neturl "net/url"
)

type Url struct {
	*neturl.URL
	IsExternal bool
}

func NewUrl(url string) (*Url, error) {
	parsed, err := neturl.Parse(url)

	return &Url{
		URL: parsed,
	}, err
}

func (u *Url) Parse(url string) (*Url, error) {
	parsed, err := u.URL.Parse(url)

	return &Url{
		URL: parsed,
	}, err
}

func (u *Url) IsEqual(url Url) bool {
	return u.Scheme == url.Scheme && u.Host == url.Host && u.Path == url.Path
}

func (u *Url) String() string {
	return u.Scheme + "://" + u.Host + u.Path
}
