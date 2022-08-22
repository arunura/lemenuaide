package main

import (
	"lemenuaide/antialexa"
	"lemenuaide/icon"

	"github.com/getlantern/systray"
)

func main() {
	onExit := func() {
		//now := time.Now()
		//ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)

}

func onReady() {
	systray.SetTemplateIcon(icon.LEMONADE, icon.LEMONADE)
	//systray.SetTitle("LeMenuAide")
	systray.SetTooltip("LeMenuAide")

	antialexa.RegisterComponent()

	// Setup Quit menu item
	mQuitOrig := systray.AddMenuItem("Quit", "Quit LeMenuAide")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

}
