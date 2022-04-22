package agent

import (
	"log"
	"net/http"
	"time"

	"github.com/digitalcircle-com-br/local-agent/lib/agent/cmd"
	"github.com/digitalcircle-com-br/local-agent/lib/agent/config"
	"github.com/digitalcircle-com-br/local-agent/lib/agent/tray"
	"github.com/digitalcircle-com-br/local-agent/lib/common"
	"github.com/gorilla/websocket"
)

func Init() error {

	err := config.Load()
	if err != nil {
		return err
	}

	go DoListen()

	tray.Run()
	return nil
}

func ListenOnce() error {

	cfg := config.Cfg

	header := http.Header{}
	header.Add("X-API-KEY", cfg.Apikey)
	header.Add("X-USER", cfg.User)
	w, _, err := websocket.DefaultDialer.Dial(cfg.Addr, header)
	if err != nil {
		return err
	}
	defer w.Close()
	
	for {
		if err != nil {
			log.Printf("Error connecting WS: %s", err.Error())
			return err
		}
		m := &common.CmdReq{}
		err = w.ReadJSON(m)
		if err != nil {
			log.Printf("Error reading req: %s", err.Error())
		} else {
			cmd.Exec(m)
		}

	}
}

func DoListen() {
	for {
		ListenOnce()
		time.Sleep(time.Second)
	}
}
