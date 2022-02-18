package chat

import (
	"log"
	"net"
)

type ChatServer struct {
	Addr            string
	Listener        net.Listener
	Clients         []*Client
	MessangerMaster chan Message
	Rooms			[]*Room
}

func NewChatServer(addr string) (ChatServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return ChatServer{}, err
	}
	return ChatServer{
		Addr:     addr,
		Listener: listener,
		MessangerMaster: make(chan Message),
		Rooms: make([]*Room, 0, 10),
	}, nil
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

func (s *ChatServer) NewRoom(creator *Client, name, password string, hidden bool) {
	m := make([]*Client, 0, 20)
	m = append(m, creator)
	r := Room{
		Name:     name,
		Password: password,
		Hidden:   hidden,
		Members:  m,
	}
	s.Rooms = append(s.Rooms, &r)
}

func (s *ChatServer) Run() {
	log.Println("Server started at ", s.Listener.Addr())
	defer s.Shutdown()

	go s.startMessanger()

	for {
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

func (s *ChatServer) handleClient(conn net.Conn) {
	c := NewClient(conn, s)
	s.Clients = append(s.Clients, &c)
	go c.Handle()
}

func (s *ChatServer) Shutdown() error {
	log.Println("Shuting down server")

	if err := s.Listener.Close(); err != nil {
		log.Println(err)
	}

	for _, c := range s.Clients {
		if err := c.Kill(); err != nil {
			log.Println(err)
		}
	}

	return nil
}
