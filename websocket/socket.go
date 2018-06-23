package websocket

import (
	"encoding/json"
	"fmt"
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

	//为了保证所有 goroutine 都可以优雅退出，需要引入 WaitGroup。陪孩子挽回，暂时先不做...
	//waitGroup *sync.WaitGroup
}

func (s *Socket) Connect(uid string, roomId int) bool {

	s.uid = uid
	s.roomId = roomId

	u := url.URL{Scheme: "ws", Host: host, Path: "/" + path}
	plog(fmt.Sprintf("%s connecting to %s\n", s.uid, u.String()))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("%s connect error:", s.uid, err)
		return false
	}

	s.wLock = new(sync.Mutex)
	s.conn = c
	s.onMessage()
	s.heartBeat()
	return true
}

func (s *Socket) onMessage() {

	go func() {
		for {
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				//log.Println("read:", err)
				plog(fmt.Sprintf("%s reade error:", s.uid, err))
				return
			}
			//log.Printf("%s recv+: %s", s.uid, message)
			plog(fmt.Sprintf("%s recv+: %s", s.uid, message))

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
				//log.Println("uid pong:", s.uid)
				plog(fmt.Sprintf("%s pong:", s.uid))

			case 5:
				//OP_SEND_SMS_REPLY
				//closeWS = true

				//统计获取到的数量

				//log.Println(p)
				s.handler.RNum += len(ptinfo)
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

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(heartHeat))
		defer ticker.Stop()

		for {
			select {
			case _ = <-ticker.C:

				buffer, err := formatHeartMessage()
				err = s.WriteMessage(buffer)
				//log.Println("uid ping:", s.uid)
				plog(fmt.Sprintf("%s ping:", s.uid))
				if err != nil {
					//log.Println("heart error", err)
					s.close()
					//TODO 通知主进程异常
					return
				}
				//case getSignal := <-s.interrupt:
				//return
			}
		}
	}()

}

func (s *Socket) Auth() bool {
	buffer, err := formatAuthMessage(s.uid, s.roomId)
	err = s.WriteMessage(buffer)
	//log.Println("write message", buffer)
	if err != nil {
		//log.Println("auth error:", err)
		plog(fmt.Sprintf("%s auth error: %s", s.uid, err))
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

	if sendInterval == 0 {
		return
	}

	go func() {

		ticker := time.NewTicker(time.Second * time.Duration(sendInterval))
		defer ticker.Stop()

		for {
			select {
			case _ = <-ticker.C:
				buffer, err := formatSendMessage(s.uid, s.roomId)
				err = s.WriteMessage(buffer)
				s.handler.PNum++
				if err != nil {
					s.close()
					plog(fmt.Sprintf("%s send error: %s", s.uid, err))
					return
				}
			}
		}
	}()
}

func (s *Socket) close() {
	if s.outPool() {
		plog(fmt.Sprintf("%s socket close", s.uid))
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
	plog(fmt.Sprintf("%s writeMessage lock", s.uid))
	s.wLock.Lock()
	defer func() {
		plog(fmt.Sprintf("%s writeMessage unlock", s.uid))
		s.wLock.Unlock()
	}()

	plog(fmt.Sprintf("%s writeMessage ==> %s", s.uid, buffer))
	err := s.conn.WriteMessage(websocket.TextMessage, buffer)
	return err
}
