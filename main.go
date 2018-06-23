package main

import (
	"benchmark/websocket"
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

	stst := h.StatisticInterval()
	ticker := time.NewTicker(time.Second * time.Duration(stst))
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			//num := h.PoolNum()
			log.Printf("uptime: %s;set client num: %d; now client num: %d; receive message num: %d; push message num:%d; in %ds", h.Uptime(), h.NumClients(), h.PoolNum(), h.ReceiveNum(), h.PushNum(), stst)
			//if num == 0 {
			//	return
			//}
		case getSignal := <-interrupt:
			log.Printf("out of %s", getSignal)
			return
		}
	}

}
