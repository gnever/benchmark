package websocket

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Socket struct {
	uid    string
	roomId int32
	conn   *websocket.Conn
}

func (s *Socket) Connect(uid string, roomId int32) {

	s.uid = uid
	s.roomId = roomId

	u := url.URL{Scheme: "ws", Host: "10.9.1.132:8090", Path: "/sub"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//defer c.Close()
	s.conn = c
	s.onMessage()
}

func (s *Socket) onMessage() {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv+: %s", message)
		}
	}()
}

func (s *Socket) Auth() {
	buffer, err := formatAuthMessage(s.uid, s.roomId)
	err = s.conn.WriteMessage(websocket.TextMessage, buffer)
	//log.Println("write message", buffer)
	if err != nil {
		log.Println("write:", err)
	}
}

func (s *Socket) Send() {
	buffer, err := formatSendMessage(s.uid, s.roomId)
	err = s.conn.WriteMessage(websocket.TextMessage, buffer)
	//log.Println("write message", buffer)
	if err != nil {
		log.Println("write:", err)
	}
}

func (s *Socket) Close() {
	s.conn.Close()
}
