package main

import (
	"log"

	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type HubOptions struct {
	MaxSavedMessage    int
	MaxReturnedMessage int
}

type Hub struct {
	Register       chan *websocket.Conn
	Unregister     chan *websocket.Conn
	GetRooms       chan *websocket.Conn
	ChangeUsername chan *Request
	JoinChat       chan *Request
	LeaveChat      chan *Request
	SendMessage    chan *Request
	OldMessages    chan *Request
	Options        *HubOptions
	connection     ConnectionStore
	user           UserStore
	room           RoomStore
	message        MessageStore
	connections    map[*websocket.Conn]*User
	users          map[string]*websocket.Conn
	mu             sync.Mutex
}

func (h *Hub) Defaults() {
	h.Options = &HubOptions{
		MaxSavedMessage:    500,
		MaxReturnedMessage: 20,
	}
	h.connection = NewInMemoryConnectionStore()
	h.user = NewInMemoryUserStore()
	h.room = NewInMemoryRoomStore()
	h.message = NewInMemoryMessageStore()
	h.connections = make(map[*websocket.Conn]*User)
	h.users = make(map[string]*websocket.Conn)
	h.mu = sync.Mutex{} // Initialize the Mutex
}

func NewHub() *Hub {
	return &Hub{
		Register:       make(chan *websocket.Conn),
		Unregister:     make(chan *websocket.Conn),
		GetRooms:       make(chan *websocket.Conn),
		ChangeUsername: make(chan *Request),
		JoinChat:       make(chan *Request),
		LeaveChat:      make(chan *Request),
		SendMessage:    make(chan *Request),
		OldMessages:    make(chan *Request),
	}
}

func (h *Hub) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		c.Locals("ClientID", uuid.String())
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Hub) Handler(ctx *fiber.Ctx) error {
	conn, err := websocket.Upgrade(ctx)
	if err != nil {
		return err
	}

	defer func() {
		h.Unregister <- conn
		if err := conn.Close(); err != nil {
			log.Printf("%#v\n", err)
		}
	}()

	h.Register <- conn

	for {
		var request Request
		if err := conn.ReadJSON(&request); err != nil {
			if e := h.error(conn, fiber.ErrBadRequest); e != nil {
				return nil
			}
			continue
		}

		// Generate request id
		id, err := uuid.NewRandom()
		if err != nil {
			if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
				return nil
			}
			continue
		}
		request.ID = id.String()

		// Set ClientID to request
		clientID, ok := ctx.Locals("ClientID").(string)
		if !ok {
			if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
				return nil
			}
			continue
		}
		request.ClientID = clientID

		// Handle incoming request based on its type
		switch request.Type {
		case GET_ROOMS:
			h.GetRooms <- conn

		case CHANGE_USERNAME:
			h.ChangeUsername <- &request
		case JOIN_CHAT:
			h.JoinChat <- &request
		case LEFT_CHAT:
			h.LeaveChat <- &request
		case SEND_MESSAGE:
			h.SendMessage <- &request
		case GET_OLD_MESSAGES:
			h.OldMessages <- &request
		// Add cases for other request types...

		default:
			if e := h.error(conn, fiber.ErrBadRequest); e != nil {
				return nil
			}
		}
	}

	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.register(conn)

		case conn := <-h.Unregister:
			h.unregister(conn)

		case conn := <-h.GetRooms:
			h.get_rooms(conn)

		case req := <-h.ChangeUsername:
			h.change_username(req)

		case req := <-h.JoinChat:
			h.join_chat(req)

		case req := <-h.LeaveChat:
			h.leave_chat(req)

		case req := <-h.SendMessage:
			h.send_message(req)

		case req := <-h.OldMessages:
			h.old_messages(req)
		}
	}
}

// Remaining methods...
