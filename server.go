package main

import (
	"github.com/gorilla/websocket"
	"net/http"
    "gqm.com/go-websocket/connection"
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
			wsConn *websocket.Conn
			data []byte
			err  error
			conn *connection.Connection
		)
		if wsConn, err = upgrader.Upgrade(writer, request, nil); err != nil {
            goto ERR
			return
		}

		if conn,err = connection.InitConnection(wsConn);err != nil{
		    goto ERR
		    return
        }
		conn.WriteMessage([]byte("connect success"))
		connection.List = append(connection.List, conn)
        for{
            if data,err = conn.ReadMessage();err != nil{
                goto ERR
            }

            if err = conn.WriteMessage(data); err!= nil{
                goto ERR
            }
        }
		ERR:
		    conn.Close()
	})

	http.ListenAndServe("0.0.0.0:8080", nil)
}

