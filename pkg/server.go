package chat

import (
	"errors"
	"log"
	"net"
)

type ChatServer struct {
	Addr            string
	Listener        net.Listener
	Clients         []*Client
	MessangerMaster chan Message
	Rooms           map[string]*Room
}

func NewChatServer(addr string) (ChatServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return ChatServer{}, err
	}
	return ChatServer{
		Addr:            addr,
		Listener:        listener,
		MessangerMaster: make(chan Message),
		Rooms:           make(map[string]*Room),
	}, nil
}

func (s *ChatServer) Run() {
	log.Println("Server started at ", s.Listener.Addr())
	defer s.disconnect()

	// start global message bus
	go s.startMessanger()

	for {
		// acccept for new clients
		conn, err := s.Listener.Accept()
		if err != nil {
			if err == net.ErrClosed {
				return
			}
			log.Println("failed accept connection:", err.Error())
			continue
		}

		s.handleClient(conn)
	}
}

func (s *ChatServer) startMessanger() {
	for {
		msg := <-s.MessangerMaster
		if msg.Client.InRoom != nil {
			for _, client := range msg.Client.InRoom.Members {
				client.FromServer <- msg
			}
		} else {
			for _, client := range s.Clients {
				if client.InRoom == nil {
					client.FromServer <- msg
				}
			}
		}

	}
}

func (s *ChatServer) NewRoom(creator *Client, name, password string, hidden bool) (*Room, error) {
	if _, found := s.Rooms[name]; found {
		return nil, errors.New("Room with this name already exists, choose another name")
	}

	m := make(map[string]*Client)
	m[creator.Addr] = creator
	r := Room{
		HelloMessage: "Welcome, welcome to " + name,
		Name:         name,
		Password:     password,
		Hidden:       hidden,
		Members:      m,
	}

	s.Rooms[name] = &r
	return &r, nil
}

func (s *ChatServer) handleClient(conn net.Conn) {
	c := NewClient(conn, s)
	s.Clients = append(s.Clients, &c)
	go c.Handle()
}

func (s *ChatServer) disconnect() {
	for _, c := range s.Clients {
		c.Connected = false
	}
}
