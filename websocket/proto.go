package websocket

import (
	"encoding/json"
)

type bodyMessage struct {
	JsessId string `json:"jsessId"`
	RoomId  int32  `json:"roomId"`
}
type proto struct {
	Ver  int32       `json:"ver"`
	Op   int16       `json:"op"`
	Seq  int16       `json:"seq"`
	Body bodyMessage `json:"body"`
}

func formatAuthMessage(uid string, roomId int32) ([]byte, error) {
	msg := bodyMessage{
		uid,
		roomId,
	}

	buff := proto{
		1,
		7,
		1,
		msg,
	}

	buffer, err := json.Marshal(buff)
	//log.Println("formatAuthMessage: ", buffer)
	return buffer, err
}

func formatSendMessage(uid string, roomId int32) ([]byte, error) {
	msg := bodyMessage{
		uid,
		roomId,
	}

	buff := proto{
		1,
		4,
		2,
		msg,
	}

	buffer, err := json.Marshal(buff)
	//log.Println("formatAuthMessage: ", buffer)
	return buffer, err
}
