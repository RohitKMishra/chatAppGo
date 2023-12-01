// room_store.go
package main

type RoomStore interface {
	AddToRoom(clientID, roomID string)
	RemoveFromRoom(clientID, roomID string)
	GetRoomMembers(roomID string) []string
}
