package shorten

import pb "github.com/aladhims/shortener/pkg/shorten/proto"

// RepositoryType is a type to identify a repository
type RepositoryType int

const (
	MongoRepository RepositoryType = 1 << iota
)

type Writer interface {
	Create(shortURL *pb.ShortURL) error
}

type Reader interface {
	GetByOrigin(origin string) *pb.ShortURL
	GetBySlug(slug string) *pb.ShortURL
	GetByBase(origin uint64) *pb.ShortURL
}

type Repository interface {
	Writer
	Reader
	Ping() error
}
