package tray

import (
	"github.com/digitalcircle-com-br/local-agent/lib/agent/config"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func onReady() {
	if len(config.Cfg.IconData)>0{
		systray.SetIcon(config.Cfg.IconData)
	}else{
		systray.SetIcon(icon.Data)
	}
	systray.SetTitle(config.Cfg.Title)
	systray.AddMenuItem(config.Cfg.User,"User")
	systray.AddMenuItem(config.Cfg.Addr,"Server")

	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func(){
		<- mQuit.ClickedCh
		systray.Quit()
	}()
	// Sets the icon of a menu item. Only available on Mac and Windows.
	
}

func onExit() {
	// clean up here
}

func Run(){
	systray.Run(onReady, onExit)
}