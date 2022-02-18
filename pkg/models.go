package chat

type Room struct {
	Name       string
	Password   string
	Hidden     bool
	Members    []*Client
	Creator    *Client
	HelloMessage string
	// MaxMembers int
}

type Message struct {
	Name	string
	Text	string
	Client  *Client
}

type CommandID int

const (
	LOGIN	CommandID	= iota
	ROOMS
	JOIN
	QUIT
	NEWROOM
)

type Command struct {
	Id	CommandID
}