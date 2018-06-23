package main

import (
	"benchmark/websocket"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	h := new(websocket.Handler)
	h.Execute()
	defer h.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			num := h.PoolNum()
			log.Println("client num:", num)
			if num == 0 {
				return
			}
		case getSignal := <-interrupt:
			fmt.Printf("out of %s", getSignal)
			return
		}
	}

}
