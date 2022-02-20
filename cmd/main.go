package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"godmodguy/chatt/pkg"
)

func main() {
	port := flag.Int("p", 6666, "port")
	flag.Parse()

	s, _ := chat.NewChatServer(fmt.Sprintf(":%d", *port))
	ctx := context.Background()
	go s.Run(ctx)


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	ctx.Done()
}