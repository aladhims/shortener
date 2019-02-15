package user

import (
	"sync"

	pb "github.com/aladhims/shortener/pkg/user/proto"
)

type Inmem struct {
	mtx       sync.RWMutex
	users     map[uint64]*pb.User
	currentId uint64
}

func NewInmemRepository() Repository {
	return &Inmem{
		users:     make(map[uint64]*pb.User),
		currentId: 1,
	}
}

func (i *Inmem) Create(user *pb.User) (uint64, error) {
	i.mtx.Lock()
	defer i.mtx.Unlock()

	i.users[i.currentId] = user
	returnedId := i.currentId

	i.currentId++

	return returnedId, nil
}

func (i *Inmem) Get(id uint64) *pb.User {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	if user, ok := i.users[id]; ok {
		return user
	}

	return nil
}

func (i *Inmem) GetByEmail(email string) *pb.User {
	i.mtx.RLock()
	defer i.mtx.RUnlock()

	for _, user := range i.users {
		if user.Email == email {
			return user
		}
	}

	return nil
}
