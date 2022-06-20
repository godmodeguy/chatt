package main

import (
	"flag"
	"fmt"
	"log"

	"godmodguy/chatt/pkg"
)

func main() {
	port := flag.Int("p", 6666, "port")
	flag.Parse()

	server, err := chat.NewChatServer(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	server.Run()
}