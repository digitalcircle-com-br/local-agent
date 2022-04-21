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

	go Do()

	tray.Run()
	return nil
}

func Once() {

	cfg := config.Cfg

	header := http.Header{}
	header.Add("X-API-KEY", cfg.Apikey)
	header.Add("X-USER", cfg.User)
	w, _, err := websocket.DefaultDialer.Dial(cfg.Addr, header)
	if err != nil {
		log.Printf("Error connecting WS: %s", err.Error())
		return
	}
	m := &common.CmdReq{}
	err = w.ReadJSON(m)
	if err != nil {
		log.Printf("Error reading req: %s", err.Error())
	}
	cmd.Exec(m)
}

func Do() {
	for {
		Once()
		time.Sleep(time.Second)
	}
}
