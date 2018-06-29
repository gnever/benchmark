package main

import (
	"benchmark/websocket"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("====START===")

	var waitgroup sync.WaitGroup

	h := new(websocket.Handler)
	h.Execute(&waitgroup)
	defer h.Close()

	stst := h.StatisticInterval()
	ticker := time.NewTicker(time.Second * time.Duration(stst))
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			//num := h.PoolNum()
			log.Printf(
				"uptime: %s;set client num: %d; now client num: %d; heartbeat num: %d; receive message num: %d; push message num:%d; in %ds",
				h.Uptime(),
				h.SetNumClients(),
				h.PoolNum(),
				h.HeartBeatNum(),
				h.ReceiveNum(),
				h.PushNum(),
				stst)
			//if num == 0 {
			//	return
			//}
		case getSignal := <-interrupt:
			log.Printf("out of %s", getSignal)
			//TODO goroutine 没有优雅的退出，可以引入waitgroup实现
			h.Close()
			waitgroup.Wait()
			log.Println("====END===")
			log.Printf(
				"uptime: %s;set client num: %d; all heartbeat num: %d; all receive message num: %d; all push message num:%d",
				h.Uptime(),
				h.SetNumClients(),
				h.AllHeartBeatNum(),
				h.AllReceiveNum(),
				h.AllPushNum())
			return
		}
	}

}
