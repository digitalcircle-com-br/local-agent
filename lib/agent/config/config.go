package config

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var Cfg = &struct {
	Addr     string
	Apikey   string
	User     string
	Cmds     map[string][]string
	Icon     string
	IconData []byte
	Title    string
	Vars     map[string]string
}{}

func Load() error {

	_, err := os.Stat("config.yaml")

	log.SetOutput(&lumberjack.Logger{
		Filename:   "agent.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	if err == nil {
		bs, err := os.ReadFile("config.yaml")
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bs, Cfg)
		if err != nil {
			return err
		}
	}

	if Cfg.Title == "" {
		Cfg.Title = " Local Agent"
	}

	if Cfg.Icon != "" {
		Cfg.IconData, err = os.ReadFile(Cfg.Icon)
		if err != nil {
			log.Printf("Error reading icon: %s", err)
		}
	}

	return nil
}
