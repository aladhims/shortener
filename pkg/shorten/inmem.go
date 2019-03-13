package shorten

import (
	"errors"
	"sync"

	pb "github.com/aladhims/shortener/pkg/shorten/proto"
)

var NotFoundError = errors.New("Not Found")

type Inmem struct {
	mtx       sync.RWMutex
	shortURLs map[uint64]*pb.ShortURL
}

func NewInmemRepository() Repository {
	return &Inmem{
		shortURLs: make(map[uint64]*pb.ShortURL),
	}
}

func (i *Inmem) Create(shortURL *pb.ShortURL) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()

	i.shortURLs[uint64(len(i.shortURLs)+1)] = shortURL

	return nil
}

func (i *Inmem) GetByOrigin(origin string) *pb.ShortURL {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	for _, url := range i.shortURLs {
		if url.Origin == origin {
			return url
		}
	}

	return nil
}

func (i *Inmem) GetBySlug(slug string) *pb.ShortURL {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	for _, url := range i.shortURLs {
		if url.Slug == slug {
			return url
		}
	}

	return nil

}

func (i *Inmem) GetByBase(base uint64) *pb.ShortURL {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	for _, url := range i.shortURLs {
		if url.Base == base {
			return url
		}
	}

	return nil
}

func (i *Inmem) Ping() error {
	return nil
}
