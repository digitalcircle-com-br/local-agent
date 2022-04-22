package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/digitalcircle-com-br/config"
	"github.com/digitalcircle-com-br/local-agent/lib/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Config = &struct {
}{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SyncMap[T any] struct {
	mx   sync.RWMutex
	data map[string]T
}

func (s *SyncMap[T]) Add(t T) string {
	s.mx.Lock()
	defer s.mx.Unlock()
	id := uuid.NewString()
	s.data[id] = t
	return id
}

func (s *SyncMap[T]) AddWID(t T, id string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[id] = t
}

func (s *SyncMap[T]) Rem(i string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, i)
}
func (s *SyncMap[T]) Get(i string) T {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.data[i]
}

func NewSync[T any]() SyncMap[T] {
	ret := SyncMap[T]{}
	ret.data = make(map[string]T)
	return ret
}

var conns = NewSync[*websocket.Conn]()
var waits = NewSync[chan []byte]()

func Run() error {

	err := config.LoadOnce(Config)
	if err != nil {
		return err
	}

	http.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		params := r.URL.Query().Get("params")
		user := r.URL.Query().Get("user")

		m := common.CmdReq{Cmd: cmd, Params: strings.Split(params, ",")}
		conn := conns.Get(user)
		conn.WriteJSON(m)

	})

	http.HandleFunc("/dowait", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		params := r.URL.Query().Get("params")
		user := r.URL.Query().Get("user")

		m := common.CmdReq{Cmd: cmd, Params: strings.Split(params, ",")}
		conn := conns.Get(user) //Find a Websocket connection to that particular user

		ch := make(chan []byte)
		id := waits.Add(ch)
		defer waits.Rem(id)
		m.ReqID = id
		err = conn.WriteJSON(m)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		select {
		case ret := <-ch:
			w.Write(ret)
			close(ch)
		case <-time.After(time.Minute * 3):
			close(ch)
			http.Error(w, "Timeout", http.StatusRequestTimeout)
		}

	})

	http.HandleFunc("/replywait", func(w http.ResponseWriter, r *http.Request) {
		// cmd := r.URL.Query().Get("cmd")
		// params := r.URL.Query().Get("params")
		// user := r.URL.Query().Get("user")
		reqid := r.URL.Query().Get("reqid")
		ch := waits.Get(reqid)

		buf := &bytes.Buffer{}
		io.Copy(buf, r.Body)
		defer r.Body.Close()

		if ch != nil {
			ch <- buf.Bytes()
		}

	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-USER")

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		conns.AddWID(conn, user) //Once a connection is made, persist that conn in a map
		ch := make(chan bool)

		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("Closing conn: %s - %v; %s", user, code, text)
			ch <- true
			conns.Rem(user)
			return nil

		})

		<-ch
		close(ch)
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	return http.ListenAndServe(":8080", nil)

}
