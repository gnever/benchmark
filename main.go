package main

import (
	"benchmark/websocket"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	s := new(websocket.Socket)
	s.Connect("asdf", 12345)
	s.Auth()
	s.Send()

	getSignal := <-interrupt
	s.Close()
	fmt.Printf("out of %s", getSignal)

}
