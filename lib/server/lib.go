package server

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/digitalcircle-com-br/local-agent/lib/common"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var mx sync.RWMutex
var conns = make(map[string]*websocket.Conn)

func AddConn(id string, c *websocket.Conn) {
	mx.Lock()
	defer mx.Unlock()
	conns[id] = c
}
func RemConn(id string) {
	mx.Lock()
	defer mx.Unlock()
	delete(conns, id)
}
func GetConn(id string) *websocket.Conn {
	mx.RLock()
	defer mx.RUnlock()
	return conns[id]
}

func Run() {
	http.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		params := r.URL.Query().Get("params")
		user := r.URL.Query().Get("user")

		m := common.CmdReq{Cmd: cmd, Params: strings.Split(params, ",")}
		conn := GetConn(user)
		conn.WriteJSON(m)

	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-USER")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		AddConn(user, conn)
		ch := make(chan bool)
		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("Closing conn: %s - %v; %s", user, code, text)
			ch <- true
			return nil
		})

		<-ch
	})

	http.ListenAndServe(":8080", nil)

}
