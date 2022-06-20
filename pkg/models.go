package chat

type Room struct {
	Name       string
	Password   string
	Hidden     bool
	Members    map[string]*Client
	Creator    *Client
	HelloMessage string
}

type Message struct {
	Name	string
	Text	string
	Client  *Client
}
