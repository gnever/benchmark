package websocket

import (
	"encoding/json"
	"log"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Socket struct {
	poolIndex int
	uid       string
	roomId    int
	conn      *websocket.Conn
	handler   *Handler
	wLock     *sync.Mutex
	//interrupt *chan os.Signal
}

func (s *Socket) Connect(uid string, roomId int) bool {

	s.uid = uid
	s.roomId = roomId

	u := url.URL{Scheme: "ws", Host: host, Path: "/" + path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("connect error:", err)
		return false
	}

	s.wLock = new(sync.Mutex)
	s.conn = c
	s.onMessage()
	go s.heartBeat()
	return true
}

func (s *Socket) onMessage() {

	go func() {
		for {
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv+: %s", message)

			var ptinfo []proto
			json.Unmarshal(message, &ptinfo)
			op := ptinfo[0].Op
			//TODO message handler

			//fmt.Println(op)

			closeWS := false

			switch op {
			case 1:
				//OP_HANDSHAKE_REPLY

			case 3:
				//OP_HEARTBEAT_REPLY
				log.Println("pong uid:", s.uid)

			case 5:
				//OP_SEND_SMS_REPLY
				//closeWS = true

				//统计获取到的数量

				//log.Println(p)
			case 8:
				//OP_AUTH_REPLY

			default:

			}

			if closeWS {
				s.close()
			}
		}
	}()
}

func (s *Socket) heartBeat() {

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:

			buffer, err := formatHeartMessage()
			err = s.WriteMessage(buffer)
			log.Println("ping uid:", s.uid)
			if err != nil {
				//log.Println("heart error", err)
				s.close()
				//TODO 通知主进程异常
				return
			}
			//case getSignal := <-s.interrupt:
			return
		}
	}

}

func (s *Socket) Auth() bool {
	buffer, err := formatAuthMessage(s.uid, s.roomId)
	err = s.WriteMessage(buffer)
	//log.Println("write message", buffer)
	if err != nil {
		log.Println("auth error:", err)
		return false
	}
	return true
}

func (s *Socket) Send() {
	//buffer, err := formatSendMessage(s.uid, s.roomId)
	//err = s.conn.WriteMessage(websocket.TextMessage, buffer)
	////log.Println("write message", buffer)
	//if err != nil {
	//	log.Println("write:", err)
	//}

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			buffer, err := formatSendMessage(s.uid, s.roomId)
			err = s.WriteMessage(buffer)
			//log.Println("write message", buffer)
			if err != nil {
				s.close()
				log.Println("send error:", err)
				return
			}
		}
	}
}

func (s *Socket) close() {
	if s.outPool() {
		log.Println("socket close:", s.uid)
		s.conn.Close()
		s.outPool()
	}
}

func (s *Socket) outPool() bool {
	if s.poolIndex == 0 {
		return false
	}

	v := reflect.ValueOf(s.handler)
	if v.IsNil() {
		return false
	}

	if _, ok := s.handler.pool.sockets[s.poolIndex]; ok {
		delete(s.handler.pool.sockets, s.poolIndex)
		return true
	}

	return false
}

func (s *Socket) WriteMessage(buffer []byte) error {
	log.Println("======lock=======")
	s.wLock.Lock()
	defer s.wLock.Unlock()

	err := s.conn.WriteMessage(websocket.TextMessage, buffer)
	return err
}
