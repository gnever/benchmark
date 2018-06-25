package websocket

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

var (
	host              string //127.0.0.1:8090
	path              string //path
	start             int
	numClients        int
	roomId            int
	heartHeat         int
	sendInterval      int
	statisticInterval int
	debug             int
)

type Handler struct {
	startTime time.Time
	pool      *pool
	RNum      int
	PNum      int
}

type pool struct {
	sockets map[int]*Socket
}

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:8090", "url")
	flag.StringVar(&path, "path", "sub", "path")
	flag.IntVar(&numClients, "m", 1, "Number of clients")
	flag.IntVar(&start, "s", 1, "start uid")
	flag.IntVar(&debug, "debug", 0, "debug")
	flag.IntVar(&roomId, "roomid", 123, "roomid")
	flag.IntVar(&heartHeat, "ht", 2, "心跳间隔")
	flag.IntVar(&sendInterval, "st", 0, "向服务端push消息间隔 单位 s, 0 为不推送")
	flag.IntVar(&statisticInterval, "stst", 5, "数据统计区间")

	flag.Parse()
	fmt.Printf("heartHeat: %ds; sendInterval: %ds;\n", heartHeat, sendInterval)
}

func (h *Handler) Execute() {
	h.startTime = time.Now()
	pool := new(pool)

	pool.sockets = make(map[int]*Socket)
	for i := start; i < start+numClients; i++ {
		s := new(Socket)
		s.poolIndex = i
		s.handler = h
		pool.sockets[i] = s

		go func() {
			s.Connect(strconv.Itoa(i), roomId)
			s.Auth()
			s.Send()
		}()
	}

	h.pool = pool

}

func (h *Handler) Close() {
	for _, s := range h.pool.sockets {
		s.close()
	}
}

func (h *Handler) PoolNum() int {
	return len(h.pool.sockets)
}

func (h *Handler) ReceiveNum() int {
	num := h.RNum
	h.RNum = 0
	return num
}

func (h *Handler) PushNum() int {
	num := h.PNum
	h.PNum = 0
	return num
}

func (h *Handler) StatisticInterval() int {
	return statisticInterval
}

func (h *Handler) NumClients() int {
	return numClients
}

func (h *Handler) Uptime() time.Duration {
	return time.Since(h.startTime)
}
