package chat

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	helloMessage = `
	Hello there!
	To quit use Ctrl-D
`

	byeMessage = `
	Bye!
`
)

type Client struct {
	Server     *ChatServer
	Username   string
	Addr       string
	Conn       net.Conn
	FromServer chan Message
	ToServer   chan Message
	InRoom     *Room
	Commander  chan Command
	Connected  bool
}

func NewClient(conn net.Conn, s *ChatServer) Client {
	return Client{
		Username:   "anonymous",
		Server:     s,
		Addr:       conn.RemoteAddr().String(),
		Conn:       conn,
		ToServer:   s.MessangerMaster,
		FromServer: make(chan Message),
		Connected:  true,
	}
}

func (c *Client) lookForMsg() {
	for {
		msg := <-c.FromServer
		var s string
		if msg.Client == c {
			continue
			// s = fmt.Sprintf("YOU]>%v\n", msg.Text)
		} else {
			roomName := "*"
			if msg.Client.InRoom != nil {
				roomName = msg.Client.InRoom.Name
			}
			s = fmt.Sprintf("[%v] %v (%v)> %v\n", roomName, msg.Name, msg.Client.Addr, msg.Text)
		}
		c.Conn.Write([]byte(s))
	}
}

func (c *Client) disconnect() {
	c.specialMessage(byeMessage)
	c.Connected = false
}

func (c *Client) Handle() {
	defer c.Conn.Close()

	log.Println("new client: ", c.Addr)
	c.specialMessage(helloMessage)

	go c.lookForMsg()

	for c.Connected {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\n\r")

		if msg == "" {
			continue
		}

		if msg == string(rune(64)) {
			c.disconnect()
			return
		}

		if strings.HasPrefix(msg, "/") {
			c.executeCommand(msg)
		} else {
			c.ToServer <- Message{
				Name:   c.Username,
				Client: c,
				Text:   msg,
			}
		}
	}
}

func (c *Client) executeCommand(cmd string) {
	commad := strings.Split(cmd, " ")
	action := strings.TrimSpace(commad[0])
	switch action {
	case "/name":
		if len(commad) == 2 {
			c.Username = commad[1]
			c.specialMessage("username changed to " + c.Username)
		} else {
			c.err("bad usage. use: /name <name>")
		}

	case "/rooms":
		rooms := make([]string, 0)
		for _, room := range c.Server.Rooms {
			if room.Hidden {
				continue
			}
			rooms = append(rooms, room.Name)
		}
		if len(rooms) == 0 {
			c.specialMessage("No available rooms")
		} else {
			msg := fmt.Sprintf(
				"%s\n%s\n%s",
				strings.Repeat("-", 20), 
				strings.Join(rooms, "\n"), 
				strings.Repeat("-", 20),
			)
			c.specialMessage(msg)
		}

	case "/join":
		password := ""
		if len(commad) < 2 || len(commad) > 3{
			c.specialMessage("bad usage. use: /join <room_name> [password]")
		}
		if len(commad) == 3 {
			password = commad[2]
		}
		err := c.joinRoom(commad[1], password)
		if err != nil {
			c.specialMessage(err.Error())
		} else {
			c.specialMessage(c.InRoom.HelloMessage)
		}

	case "/quit":
		if c.InRoom != nil {
			c.specialMessage("exiting...")
			c.quitRoom()
		} else {
			c.disconnect()
		}

	case "/newroom":
		if len(commad) < 2 || strings.Contains(commad[1], ":") {
			c.err("bad usage. use: /newroom <room_name> p:[password] [-h (hidden)]")
			return
		}
		name := commad[1]
		password := ""
		hidden := false
		for _, arg := range commad[1:] {
			if strings.HasPrefix(arg, "p:") {
				password = arg[2:]
			}
			if arg == "-h" {
				hidden = true
			}
		}
		room, err := c.Server.NewRoom(c, name, password, hidden)
		c.InRoom = room
		if err != nil {
			c.specialMessage(err.Error())
		} else {
			c.specialMessage("Created")
		}

	default:
		c.err(fmt.Sprintf("unknown command: %v", action))
	}
}

func (c *Client) joinRoom(name, password string) error {
	room, found := c.Server.Rooms[name]
	if !found {
		return errors.New("no such room")
	}

	if room.Password != password {
		return errors.New("incorrect password")
	}

	c.InRoom = room
	c.Server.Rooms[name].Members[c.Addr] = c
	return nil
} 

func (c *Client) quitRoom() {
	mems := c.Server.Rooms[c.InRoom.Name].Members
	delete(mems, c.Addr)
	c.InRoom = nil
}

func (c *Client) err(e string) {
	c.Conn.Write([]byte("ERR " + e + "\n"))
}

func (c *Client) specialMessage(s string) {
	c.Conn.Write([]byte(s + "\n"))
}

func (c *Client) Kill() error {
	return c.Conn.Close()
}
