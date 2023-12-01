// in_memory_room_store.go
package main

import (
	"sync"
)

type InMemoryRoomStore struct {
	sync.Mutex
	rooms map[string]map[string]struct{} // roomID -> set of clientIDs
}

var _ RoomStore = (*InMemoryRoomStore)(nil)

func NewInMemoryRoomStore() *InMemoryRoomStore {
	return &InMemoryRoomStore{
		rooms: make(map[string]map[string]struct{}),
	}
}

func (s *InMemoryRoomStore) AddToRoom(clientID, roomID string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.rooms[roomID]; !ok {
		s.rooms[roomID] = make(map[string]struct{})
	}

	s.rooms[roomID][clientID] = struct{}{}
}

func (s *InMemoryRoomStore) RemoveFromRoom(clientID, roomID string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.rooms[roomID]; ok {
		delete(s.rooms[roomID], clientID)

		if len(s.rooms[roomID]) == 0 {
			delete(s.rooms, roomID)
		}
	}
}

func (s *InMemoryRoomStore) GetRoomMembers(roomID string) []string {
	s.Lock()
	defer s.Unlock()

	members := make([]string, 0, len(s.rooms[roomID]))
	for clientID := range s.rooms[roomID] {
		members = append(members, clientID)
	}

	return members
}
