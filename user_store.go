// user_store.go
package main

import (
	"sync"

	"github.com/RohitKMishra/chatAppGo/models"
)

type UserStore interface {
	Store(userID string, user models.User)
	Load(userID string) (user models.User, ok bool)
	Delete(userID string)
}

type InMemoryUserStore struct {
	sync.Mutex
	users map[string]models.User
}

var _ UserStore = (*InMemoryUserStore)(nil)

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: map[string]models.User{},
	}
}

func (s *InMemoryUserStore) Store(userID string, user models.User) {
	s.Lock()
	s.users[userID] = user
	s.Unlock()
}

func (s *InMemoryUserStore) Load(userID string) (user models.User, ok bool) {
	s.Lock()
	user, ok = s.users[userID]
	s.Unlock()
	return user, ok
}

func (s *InMemoryUserStore) Delete(userID string) {
	s.Lock()
	delete(s.users, userID)
	s.Unlock()
}
