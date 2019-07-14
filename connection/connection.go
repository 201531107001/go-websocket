package connection

import (
    "github.com/gorilla/websocket"
    "sync"
    "errors"
    "fmt"
)

var List = make([]*Connection,0)
var id = 1

type Connection struct {
    wsConn *websocket.Conn
    inChan chan []byte
    outChan chan []byte
    closeChan chan byte

    mu sync.Mutex
    isClosed bool

    connId int
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection,err error) {
    conn = &Connection{
        wsConn:wsConn,
        inChan:make(chan []byte,1000),
        outChan:make(chan []byte,1000),
        closeChan:make(chan byte,1),
        isClosed:false,
        connId: id,
    }
    id++
    go conn.readLoop()
    go conn.writeLoop()
    return
}

func (conn *Connection)ReadMessage() (data []byte,err error) {
    select {
    case data = <- conn.inChan:
    case <- conn.closeChan:
        err = errors.New("conn is closed")
    }
    return
}

func (conn *Connection)WriteMessage(data []byte) (err error) {
    select {
    case conn.outChan <- data:
    fmt.Println("toAll")
    case <- conn.closeChan:
        err = errors.New("conn is closed")
    }
    return
}

func (conn Connection)Close() {
    conn.WriteMessage([]byte("connect closed"))
    conn.wsConn.Close()
    conn.mu.Lock()
    if !conn.isClosed {
        close(conn.closeChan)
        conn.isClosed = true
    }
    conn.mu.Unlock()
}

func (conn *Connection)readLoop() {
    var (
       date []byte
       err error
    )
    for  {
        _,date,err = conn.wsConn.ReadMessage()
        if err != nil {
            goto ERR
        }
        select {
            case conn.inChan <- date:
                toAll(conn,date)
        case <- conn.closeChan:
            goto ERR
        }
    }
ERR:
    conn.Close()
}

func (conn Connection)writeLoop() {
    var (
        date []byte
        err error
    )
    for  {
        select {
        case date = <- conn.outChan:
        case <- conn.closeChan:
            goto ERR
        }
        err = conn.wsConn.WriteMessage(websocket.TextMessage,date)
        if err != nil {
            goto ERR
        }
    }
ERR:
    conn.Close()
}

func toAll(conn *Connection,mess []byte)  {
    for _,ot := range List{
       if ot.connId != conn.connId{
          ot.WriteMessage(mess)
       }
    }
}