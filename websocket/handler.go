package websocket

import (
	"flag"
	"fmt"
	"strconv"
)

var (
	host       string //127.0.0.1:8090
	path       string //path
	start      int
	numClients int
	roomId     int
	heartHeat  int
)

type Handler struct {
	pool *pool
}

type pool struct {
	sockets map[int]*Socket
}

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:8090", "url")
	flag.StringVar(&path, "path", "sub", "path")
	flag.IntVar(&numClients, "m", 1, "Number of clients")
	flag.IntVar(&start, "s", 1, "start uid")
	flag.IntVar(&roomId, "roomid", 123, "roomid")
	flag.IntVar(&heartHeat, "ht", 2, "心跳间隔")

	flag.Parse()

	fmt.Println(host)
	fmt.Println(path)
	fmt.Println(numClients)
}

func (h *Handler) Execute() {
	pool := new(pool)

	pool.sockets = make(map[int]*Socket)
	for i := start; i < start+numClients; i++ {
		s := new(Socket)
		s.poolIndex = i
		s.handler = h
		pool.sockets[i] = s

		s.Connect(strconv.Itoa(i), roomId)
		s.Auth()
		//s.Send()
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
