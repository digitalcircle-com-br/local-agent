package agent

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/digitalcircle-com-br/local-agent/lib/agent/tray"
	"github.com/digitalcircle-com-br/local-agent/lib/common"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr   string
	Apikey string
	User   string
	Cmds   map[string]string
	Icon   string
}

var cfg *Config

func Init() error {
	cfg = &Config{}
	_, err := os.Stat("config.yaml")
	if err == nil {
		bs, err := os.ReadFile("config.yaml")
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bs, cfg)
		if err != nil {
			return err
		}
	}
	go Do()

	tray.Run()
	// for {
	// 	time.Sleep(time.Minute)
	// }
	return nil
}

func Once() {
	header := http.Header{}
	header.Add("X-API-KEY", cfg.Apikey)
	header.Add("X-USER", cfg.User)
	w, _, err := websocket.DefaultDialer.Dial(cfg.Addr, header)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	m := common.CmdReq{}
	err = w.ReadJSON(&m)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("%v", m)
}

func Do() {
	for {
		Once()
		time.Sleep(time.Second)
	}
}
