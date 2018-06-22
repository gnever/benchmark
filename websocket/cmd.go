package websocket

import (
	"flag"
	"fmt"
)

var (
	host       string //127.0.0.1:8090
	path       string //path
	numClients int
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1:8090", "url")
	flag.StringVar(&path, "p", "sub", "path")
	flag.IntVar(&numClients, "m", 1, "Number of clients.")
	flag.Parse()

	fmt.Println(host)
	fmt.Println(path)
	fmt.Println(numClients)
}
