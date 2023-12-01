// in_memory_message_store.go
package main

import (
	"sync"
)

type InMemoryMessageStore struct {
	sync.Mutex
	messages map[string][]Message // roomID -> slice of messages
}

var _ MessageStore = (*InMemoryMessageStore)(nil)

func NewInMemoryMessageStore() *InMemoryMessageStore {
	return &InMemoryMessageStore{
		messages: make(map[string][]Message),
	}
}

func (s *InMemoryMessageStore) Append(roomID string, message Message) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.messages[roomID]; !ok {
		s.messages[roomID] = make([]Message, 0)
	}

	s.messages[roomID] = append(s.messages[roomID], message)
}

func (s *InMemoryMessageStore) GetLastN(roomID string, n int, beforeID ...string) []Message {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.messages[roomID]; !ok {
		return nil
	}

	var startIndex int
	if len(beforeID) > 0 {
		for i, msg := range s.messages[roomID] {
			if msg.ID == beforeID[0] {
				startIndex = i + 1
				break
			}
		}
	}

	length := len(s.messages[roomID])
	if length <= startIndex {
		return nil
	}

	endIndex := length
	if length-startIndex > n {
		endIndex = startIndex + n
	}

	return s.messages[roomID][startIndex:endIndex]
}

func (s *InMemoryMessageStore) Count(roomID string) int {
	s.Lock()
	defer s.Unlock()

	return len(s.messages[roomID])
}
