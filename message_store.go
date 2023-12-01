// message_store.go
package main

type MessageStore interface {
	Append(roomID string, message Message)
	GetLastN(roomID string, n int, beforeID ...string) []Message
	Count(roomID string) int
}
