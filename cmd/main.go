package main

import (
	"os"
	"os/signal"
	"syscall"

	"godmodguy/chatt/pkg"
)

func main() {
	s, _ := chat.NewChatServer(":5599")
	go s.Run()



	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	// s.Shutdown()
}