package main

import "github.com/google/uuid"

const (
	GET_ROOMS        = "get_rooms"
	CHANGE_USERNAME  = "change_username"
	JOIN_CHAT        = "join_chat"
	LEFT_CHAT        = "left_chat"
	SEND_MESSAGE     = "send_message"
	GET_OLD_MESSAGES = "get_old_messages"
)

type Request struct {
	ID        string `json:"id,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	RoomID    string `json:"room_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type Message struct {
	ID        string `json:"id,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	RoomID    string `json:"room_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type Room struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
}
