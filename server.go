package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		var (
			conn *websocket.Conn
			date []byte
			err  error
		)
		if conn, err = upgrader.Upgrade(writer, request, nil); err != nil {
			return
		}
        conn.WriteMessage(websocket.TextMessage, []byte("连接成功"))
		for true {
			if _, date, err = conn.ReadMessage(); err != nil {
                conn.WriteMessage(websocket.TextMessage, []byte("断开连接"))
				conn.Close()
				fmt.Println("断开连接", err)
				return
			}
			fmt.Println(string(date))
			conn.WriteMessage(websocket.TextMessage, date)
		}
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}

