package user

import pb "github.com/aladhims/shortener/pkg/user/proto"

// RepositoryType is a type to identify a repository
type RepositoryType int

const (
	PostgreRepository RepositoryType = 1 << iota
)

type Writer interface {
	Create(user *pb.User) (uint64, error)
}

type Reader interface {
	Get(id uint64) *pb.User
	GetByEmail(email string) *pb.User
}

type Repository interface {
	Writer
	Reader
	Ping() error
}
