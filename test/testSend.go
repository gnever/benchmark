// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
    "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type bodyMessage struct {
    JsessId string `json:"jsessId"`
    RoomId int32 `json:"roomId"`
}
type proto struct {
    Ver int32 `json:"ver"`
    Op int16 `json:"op"`
    Seq int16 `json:"seq"`
    Body bodyMessage `json:"body"`
}
//type Proto struct {
//    PackLen   int32  // package length
//    HeaderLen int16  // header length
//    Ver       int16  // protocol version
//    Operation int32  // operation for request
//    SeqId     int32  // sequence number chosen by client
//    Body      []byte // body
//}

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8090", Path: "/sub"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()



    msg := bodyMessage {
            "9bn4t1mr8dhlceuhkb941tp8u1",
            32793,
        }

    buff := proto {
        1,
        7,
        1,
        msg,
    }


    buffer, err := json.Marshal(buff)
    fmt.Printf("%s\n", buffer)
    //buffer := "{ver: 1, op: 7, seq: 1, body: {jsessId: "9bn4t1mr8dhlceuhkb941tp8u1", roomId: 32793}}";
    //buffer := "{"ver":1,"op":3,"seq":1,"body":[{"jsessId":"9bn4t1mr8dhlceuhkb941tp8u1","roomId":32793}]}"


    fmt.Println("write buffer")
    err = c.WriteMessage(websocket.TextMessage, buffer)
    log.Println("write message", buffer)
    if err != nil {
        log.Println("write:", err)
        return
    }


	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:

			fmt.Println("write buffer")

            msg := bodyMessage {
                    "9bn4t1mr8dhlceuhkb941tp8u1",
                    32793,
            }

            buff := proto {
                1,
                4,
                2,
                msg,
            }


            buffer, err := json.Marshal(buff)
            fmt.Printf("%s\n", buffer)
            //buffer := "{ver: 1, op: 7, seq: 1, body: {jsessId: "9bn4t1mr8dhlceuhkb941tp8u1", roomId: 32793}}";
            //buffer := "{"ver":1,"op":3,"seq":1,"body":[{"jsessId":"9bn4t1mr8dhlceuhkb941tp8u1","roomId":32793}]}"


            fmt.Println("write buffer")
            err = c.WriteMessage(websocket.TextMessage, buffer)
            log.Println("write message", buffer)
            if err != nil {
                log.Println("write:", err)
                return
            }



        case <-interrupt:
            log.Println("interrupt")

            // Cleanly close the connection by sending a close message and then
            // waiting (with timeout) for the server to close the connection.
            err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("write close:", err)
                return
            }
            select {
            case <-done:
            case <-time.After(time.Second):
            }
            return
        }
    }



}
